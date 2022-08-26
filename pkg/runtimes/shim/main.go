package main

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/containerd/containerd/sys/reaper"
	"github.com/mattn/go-colorable"
	"github.com/moby/sys/mountinfo"
	"github.com/opencontainers/go-digest"
	ispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/opencontainers/umoci/oci/cas"
	"github.com/opencontainers/umoci/oci/casext"
	"github.com/opencontainers/umoci/oci/layer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sys/unix"
)

// NB: change this to Debug if you're troubleshooting the shim
const LogLevel = zapcore.ErrorLevel

type Command struct {
	Args  []string `json:"args"`
	Stdin []byte   `json:"stdin"`
	Env   []string `json:"env"`
	Dir   *string  `json:"dir"`
}

var stdoutPath string

func init() {
	stdoutPath = os.Getenv("_BASS_OUTPUT")
	os.Unsetenv("_BASS_OUTPUT")
}

var cmds = map[string]func([]string) error{
	"run":        run,
	"unpack":     unpack,
	"get-config": getConfig,
	"check":      check,
}

var cmdArg string

func init() {
	var cmdOpts []string
	for k := range cmds {
		cmdOpts = append(cmdOpts, k)
	}

	sort.Strings(cmdOpts)

	cmdArg = strings.Join(cmdOpts, "|")
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "usage: %s <%s>\n", os.Args[0], cmdArg)
		os.Exit(1)
	}

	cmd, args := os.Args[1], os.Args[2:]

	f, found := cmds[cmd]
	if !found {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		fmt.Fprintf(os.Stderr, "usage: %s <%s>\n", os.Args[0], cmdArg)
		os.Exit(1)
		return
	}

	err := f(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: run <cmd.json>")
	}

	err := installCert()
	if err != nil {
		return fmt.Errorf("install bass CA: %w", err)
	}

	err = reap()
	if err != nil {
		return fmt.Errorf("reap: %w", err)
	}

	cmdPath := args[0]

	cmdPayload, err := os.ReadFile(cmdPath)
	if err != nil {
		return fmt.Errorf("read cmd: %w", err)
	}

	var cmd Command
	err = json.Unmarshal(cmdPayload, &cmd)
	if err != nil {
		return fmt.Errorf("unmarshal cmd: %w", err)
	}

	err = os.Remove(cmdPath)
	if err != nil {
		return fmt.Errorf("burn after reading: %w", err)
	}

	var stdout io.Writer = os.Stdout
	if stdoutPath != "" {
		response, err := os.Create(stdoutPath)
		if err != nil {
			return fmt.Errorf("create output error: %w", err)
		}

		defer response.Close()

		stdout = response
	}

	for _, e := range cmd.Env {
		segs := strings.SplitN(e, "=", 2)
		if len(segs) != 2 {
			return fmt.Errorf("malformed env: %s", e)
		}

		os.Setenv(segs[0], segs[1])
	}

	bin := cmd.Args[0]
	argv := cmd.Args[1:]
	execCmd := exec.Command(bin, argv...)
	if cmd.Dir != nil {
		execCmd.Dir = *cmd.Dir
	}
	execCmd.Stdin = bytes.NewBuffer(cmd.Stdin)
	execCmd.Stdout = stdout
	execCmd.Stderr = os.Stderr

	ch, err := reaper.Default.Start(execCmd)
	if err != nil {
		return fmt.Errorf("start: %w", err)
	}

	status, err := reaper.Default.Wait(execCmd, ch)
	if err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	if status != 0 {
		// propagate exit status
		os.Exit(status)
		return nil
	}

	err = normalizeTimes(".")
	if err != nil {
		return fmt.Errorf("failed to normalize timestamps: %w", err)
	}

	return nil
}

func reap() error {
	logger := StdLogger(LogLevel)

	reaper.SetSubreaper(1)

	children := make(chan os.Signal, 32)
	signal.Notify(children, syscall.SIGCHLD)

	go func() {
		for range children {
			err := reaper.Reap()
			if err != nil {
				logger.Warn("failed to reap", zap.Error(err))
			}
		}
	}()

	return nil
}

func getConfig(args []string) error {
	ctx := context.Background()

	if len(args) != 3 {
		return fmt.Errorf("usage: get-config image.tar tag dest/")
	}

	archiveSrc := args[0]
	fromName := args[1]
	configDst := args[2]

	layout, err := openTar(archiveSrc)
	if err != nil {
		return fmt.Errorf("create layout: %w", err)
	}

	defer layout.Close()

	ext := casext.NewEngine(layout)

	mspec, err := loadManifest(ctx, ext, fromName)
	if err != nil {
		return err
	}

	config, err := ext.FromDescriptor(ctx, mspec.Config)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if config.Descriptor.MediaType != ispec.MediaTypeImageConfig {
		return fmt.Errorf("bad config media type: %s", config.Descriptor.MediaType)
	}

	ispec := config.Data.(ispec.Image)

	configPath := filepath.Join(configDst, "config.json")

	configFile, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("create config.json: %w", err)
	}

	defer configFile.Close()

	err = json.NewEncoder(configFile).Encode(ispec.Config)
	if err != nil {
		return fmt.Errorf("encode image config: %w", err)
	}

	return nil
}

func unpack(args []string) error {
	ctx := context.Background()

	if len(args) != 3 {
		return fmt.Errorf("usage: unpack <image.tar> <tag> <dest/>")
	}

	archiveSrc := args[0]
	fromName := args[1]
	rootfsPath := args[2]

	layout, err := openTar(archiveSrc)
	if err != nil {
		return fmt.Errorf("create layout: %w", err)
	}

	defer layout.Close()

	ext := casext.NewEngine(layout)

	mspec, err := loadManifest(ctx, ext, fromName)
	if err != nil {
		return err
	}

	err = layer.UnpackRootfs(context.Background(), ext, rootfsPath, mspec, &layer.UnpackOptions{})
	if err != nil {
		return fmt.Errorf("unpack rootfs: %w", err)
	}

	return nil
}

func loadManifest(ctx context.Context, ext casext.Engine, name string) (ispec.Manifest, error) {
	descPaths, err := ext.ResolveReference(context.Background(), name)
	if err != nil {
		return ispec.Manifest{}, fmt.Errorf("resolve ref: %w", err)
	}

	if len(descPaths) == 0 {
		return ispec.Manifest{}, fmt.Errorf("tag not found: %s", name)
	}

	if len(descPaths) != 1 {
		return ispec.Manifest{}, fmt.Errorf("ambiguous tag?: %s (%d paths returned)", name, len(descPaths))
	}

	manifest, err := ext.FromDescriptor(ctx, descPaths[0].Descriptor())
	if err != nil {
		return ispec.Manifest{}, fmt.Errorf("load manifest: %w", err)
	}

	if manifest.Descriptor.MediaType != ispec.MediaTypeImageManifest {
		return ispec.Manifest{}, fmt.Errorf("bad manifest media type: %s", manifest.Descriptor.MediaType)
	}

	return manifest.Data.(ispec.Manifest), nil
}

func openTar(tarPath string) (cas.Engine, error) {
	archive, err := os.Open(tarPath)
	if err != nil {
		return nil, err
	}

	return &tarEngine{archive}, nil
}

// tarEngine implements a read-only cas.Engine backed by a .tar archive.
type tarEngine struct {
	archive *os.File
}

func (engine *tarEngine) PutBlob(ctx context.Context, reader io.Reader) (digest.Digest, int64, error) {
	return "", 0, fmt.Errorf("PutBlob: %w", cas.ErrNotImplemented)
}

func (engine *tarEngine) GetBlob(ctx context.Context, dig digest.Digest) (io.ReadCloser, error) {
	r, err := engine.open(path.Join("blobs", dig.Algorithm().String(), dig.Encoded()))
	if err != nil {
		return nil, err
	}

	return io.NopCloser(r), nil
}

func (engine *tarEngine) StatBlob(ctx context.Context, dig digest.Digest) (bool, error) {
	_, err := engine.open(path.Join("blobs", dig.Algorithm().String(), dig.Encoded()))
	if err != nil {
		if errors.Is(err, cas.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (engine *tarEngine) PutIndex(ctx context.Context, index ispec.Index) error {
	return fmt.Errorf("PutIndex: %w", cas.ErrNotImplemented)
}

func (engine *tarEngine) GetIndex(ctx context.Context) (ispec.Index, error) {
	var idx ispec.Index
	r, err := engine.open("index.json")
	if err != nil {
		return ispec.Index{}, err
	}

	err = json.NewDecoder(r).Decode(&idx)
	if err != nil {
		return ispec.Index{}, err
	}

	return idx, nil
}

func (engine *tarEngine) open(p string) (io.Reader, error) {
	_, err := engine.archive.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("seek: %w", err)
	}

	tr := tar.NewReader(engine.archive)

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("open %s: %w", p, cas.ErrNotExist)
			}

			return nil, err
		}

		if path.Clean(hdr.Name) == p {
			return tr, nil
		}
	}
}

func (engine *tarEngine) DeleteBlob(ctx context.Context, digest digest.Digest) (err error) {
	return fmt.Errorf("DeleteBlob: %w", cas.ErrNotImplemented)
}

func (engine *tarEngine) ListBlobs(ctx context.Context) ([]digest.Digest, error) {
	_, err := engine.archive.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("seek: %w", err)
	}

	tr := tar.NewReader(engine.archive)

	var digs []digest.Digest
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("next: %w", err)
		}

		if strings.HasPrefix(path.Clean(hdr.Name), "blobs/") {
			dir, encoded := path.Split(hdr.Name)
			_, alg := path.Split(dir)
			digs = append(digs, digest.NewDigestFromEncoded(digest.Algorithm(alg), encoded))
		}
	}

	return digs, nil
}

func (engine *tarEngine) Clean(ctx context.Context) error { return nil }

func (engine *tarEngine) Close() error {
	return engine.archive.Close()
}

var epoch = time.Date(1985, 10, 26, 8, 15, 0, 0, time.UTC)

func normalizeTimes(root string) error {
	logger := StdLogger(LogLevel)

	skipped := 0
	unchanged := 0
	changed := 0
	start := time.Now()
	tspec := unix.NsecToTimespec(epoch.UnixNano())
	targetTime := []unix.Timespec{tspec, tspec}
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if path != root && info.IsDir() {
			mp, err := mountinfo.Mounted(path)
			if err != nil {
				return fmt.Errorf("check mounted: %w", err)
			}

			if mp {
				logger.Debug("skipping mountpoint", zap.String("path", path))
				skipped++
				return fs.SkipDir
			}
		}

		if info.ModTime().Equal(epoch) {
			unchanged++
			return nil
		}

		changed++

		logger.Debug("chtimes",
			zap.String("path", path),
			zap.Time("from", info.ModTime()),
			zap.Time("to", epoch))

		err = unix.UtimesNanoAt(unix.AT_FDCWD, path, targetTime, unix.AT_SYMLINK_NOFOLLOW)
		if err != nil {
			return fmt.Errorf("chtimes: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("times normalized",
		zap.Duration("took", time.Since(start)),
		zap.Int("changed", changed),
		zap.Int("unchanged", unchanged),
		zap.Int("skipped", skipped),
	)

	return nil
}

func check(args []string) error {
	logger := StdLogger(LogLevel)

	if len(args) == 0 {
		return fmt.Errorf("usage: check <host> name:port...")
	}

	host, ports := args[0], args[1:]

	for _, nameAndPort := range ports {
		name, port, ok := strings.Cut(nameAndPort, ":")
		if !ok {
			return fmt.Errorf("port must be in form name:number: %s", nameAndPort)
		}

		logger := logger.With(zap.String("name", name), zap.String("port", port))

		logger.Debug("polling for port")

		pollAddr := net.JoinHostPort(host, port)

		err := pollForPort(logger, pollAddr)
		if err != nil {
			return fmt.Errorf("poll %s: %w", name, err)
		}

		logger.Info("port is up")
	}

	return nil
}

func pollForPort(logger *zap.Logger, addr string) error {
	retry := backoff.NewExponentialBackOff()
	retry.InitialInterval = 100 * time.Millisecond

	dialer := net.Dialer{
		Timeout: time.Second,
	}

	return backoff.Retry(func() error {
		conn, err := dialer.Dial("tcp", addr)
		if err != nil {
			logger.Debug("failed to dial", zap.Duration("elapsed", retry.GetElapsedTime()), zap.Error(err))
			return err
		}

		_ = conn.Close()

		return nil
	}, retry)
}

// yoinked from pkg/bass/log.go, avoiding too many dependencies
func LoggerTo(w io.Writer, level zapcore.LevelEnabler) *zap.Logger {
	zapcfg := zap.NewDevelopmentEncoderConfig()
	zapcfg.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	zapcfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("15:04:05.000"))
	}

	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcfg),
		zapcore.AddSync(w),
		level,
	))
}

func StdLogger(level zapcore.LevelEnabler) *zap.Logger {
	return LoggerTo(colorable.NewColorableStderr(), level)
}
