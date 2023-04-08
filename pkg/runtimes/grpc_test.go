package runtimes_test

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vito/bass/pkg/bass"
	"github.com/vito/bass/pkg/proto"
	"github.com/vito/bass/pkg/runtimes"
	"github.com/vito/bass/pkg/runtimes/util/buildkitd"
	"github.com/vito/is"
	"google.golang.org/grpc"
)

func TestGRPCRuntime(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	const testInst = "bass-buildkitd-test"

	buildkitd.Remove(context.Background(), testInst)

	config := bass.Bindings{
		"debug":        bass.Bool(true),
		"installation": bass.String(testInst),
	}

	if dir, ok := os.LookupEnv("BASS_TLS_DEPOT"); ok && dir != "" {
		config["certs_dir"] = bass.String(dir)
	}

	sockPath := filepath.Join(t.TempDir(), "sock")
	listener, err := net.Listen("unix", sockPath)
	require.NoError(t, err)

	defer listener.Close()

	ctx := context.Background()

	pool, err := runtimes.NewPool(ctx, &bass.Config{
		Runtimes: []bass.RuntimeConfig{
			{
				Platform: bass.LinuxPlatform,
				Runtime:  runtimes.BuildkitName,
			},
		},
	})
	is.New(t).NoErr(err)

	srv := grpc.NewServer()
	proto.RegisterRuntimeServer(srv, &runtimes.Server{
		Context: ctx,
		Runtime: pool.Runtimes[0].Runtime,
	})

	go func() {
		if err := srv.Serve(listener); err != nil && !errors.Is(err, net.ErrClosed) {
			panic(err)
		}
	}()

	runtimes.Suite(t, bass.RuntimeConfig{
		Platform: bass.LinuxPlatform,
		Runtime:  runtimes.GRPCName,
		Config: bass.Bindings{
			"target": bass.String("unix://" + sockPath),
		}.Scope(),
	})
}
