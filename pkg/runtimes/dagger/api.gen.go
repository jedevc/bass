// Code generated by dagger. DO NOT EDIT.

package dagger

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"go.dagger.io/dagger/sdk/go/dagger/querybuilder"
)

// New returns a new API query object
func New(c graphql.Client) *Query {
	return &Query{
		q: querybuilder.Query(),
		c: c,
	}
}

// A global cache volume identifier
type CacheID string

// The address (also known as "ref") of a container published as an OCI image.
//
// Examples:
//   - "alpine"
//   - "index.docker.io/alpine"
//   - "index.docker.io/alpine:latest"
//   - "index.docker.io/alpine:latest@sha256deadbeefdeadbeefdeadbeef"
type ContainerAddress string

// A unique container identifier. Null designates an empty container (scratch).
type ContainerID string

// A content-addressed directory identifier
type DirectoryID string

type FSID string

type FileID string

// An identifier for a directory on the host
type HostDirectoryID string

// A unique identifier for a secret
type SecretID string

type CacheMountInput struct {
	// Cache mount name
	Name string `json:"name"`

	// path at which the cache will be mounted
	Path string `json:"path"`

	// Cache mount sharing mode (TODO: switch to enum)
	SharingMode string `json:"sharingMode"`
}

type ExecEnvInput struct {
	// Env var name
	Name string `json:"name"`

	// Env var value
	Value string `json:"value"`
}

type ExecInput struct {
	// Command to execute
	// Example: ["echo", "hello, world!"]
	Args []string `json:"args"`

	// Cached mounts
	CacheMounts []CacheMountInput `json:"cacheMounts"`

	// Env vars
	Env []ExecEnvInput `json:"env"`

	// Filesystem mounts
	Mounts []MountInput `json:"mounts"`

	// Secret env vars
	SecretEnv []ExecSecretEnvInput `json:"secretEnv"`

	// Include the host's ssh agent socket in the exec at the provided path
	SSHAuthSock string `json:"sshAuthSock"`

	// Working directory
	Workdir string `json:"workdir"`
}

// Additional options for executing a command
type ExecOpts struct {
	// Optionally redirect the command's standard error to a file in the container.
	// Null means discard output.
	RedirectStderr string `json:"redirectStderr"`

	// Optionally redirect the command's standard output to a file in the container.
	// Null means discard output.
	RedirectStdout string `json:"redirectStdout"`

	// Optionally write to the command's standard input
	//
	// - Null means don't touch stdin (no redirection)
	// - Empty string means inject zero bytes to stdin, then send EOF
	Stdin string `json:"stdin"`
}

type ExecSecretEnvInput struct {
	// Secret env var value
	ID SecretID `json:"id"`

	// Env var name
	Name string `json:"name"`
}

type MountInput struct {
	// filesystem to mount
	FS FSID `json:"fs"`

	// path at which the filesystem will be mounted
	Path string `json:"path"`
}

// A directory whose contents persist across runs
type CacheVolume struct {
	q *querybuilder.Selection
	c graphql.Client
}

func (r *CacheVolume) ID(ctx context.Context) (CacheID, error) {
	q := r.q.Select("id")

	var response CacheID
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// An OCI-compatible container, also known as a docker container
type Container struct {
	q *querybuilder.Selection
	c graphql.Client
}

// Default arguments for future commands
func (r *Container) DefaultArgs(ctx context.Context) ([]string, error) {
	q := r.q.Select("defaultArgs")

	var response []string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// Retrieve a directory at the given path. Mounts are included.
func (r *Container) Directory(path string) *Directory {
	q := r.q.Select("directory")
	q = q.Arg("path", path)

	return &Directory{
		q: q,
		c: r.c,
	}
}

// Entrypoint to be prepended to the arguments of all commands
func (r *Container) Entrypoint(ctx context.Context) ([]string, error) {
	q := r.q.Select("entrypoint")

	var response []string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// ContainerExecOpts contains options for Container.Exec
type ContainerExecOpts struct {
	Args []string

	Opts ExecOpts
}

// This container after executing the specified command inside it
func (r *Container) Exec(opts ...ContainerExecOpts) *Container {
	q := r.q.Select("exec")
	// `args` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Args) {
			q = q.Arg("args", opts[i].Args)
			break
		}
	}
	// `opts` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Opts) {
			q = q.Arg("opts", opts[i].Opts)
			break
		}
	}

	return &Container{
		q: q,
		c: r.c,
	}
}

// Exit code of the last executed command. Zero means success.
// Null if no command has been executed.
func (r *Container) ExitCode(ctx context.Context) (int, error) {
	q := r.q.Select("exitCode")

	var response int
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// Retrieve a file at the given path. Mounts are included.
func (r *Container) File(path string) *File {
	q := r.q.Select("file")
	q = q.Arg("path", path)

	return &File{
		q: q,
		c: r.c,
	}
}

// Initialize this container from the base image published at the given address
func (r *Container) From(address ContainerAddress) *Container {
	q := r.q.Select("from")
	q = q.Arg("address", address)

	return &Container{
		q: q,
		c: r.c,
	}
}

// A unique identifier for this container
func (r *Container) ID(ctx context.Context) (ContainerID, error) {
	q := r.q.Select("id")

	var response ContainerID
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// List of paths where a directory is mounted
func (r *Container) Mounts(ctx context.Context) ([]string, error) {
	q := r.q.Select("mounts")

	var response []string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// Publish this container as a new image
func (r *Container) Publish(ctx context.Context, address ContainerAddress) (ContainerAddress, error) {
	q := r.q.Select("publish")
	q = q.Arg("address", address)

	var response ContainerAddress
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// This container's root filesystem. Mounts are not included.
func (r *Container) Rootfs() *Directory {
	q := r.q.Select("rootfs")

	return &Directory{
		q: q,
		c: r.c,
	}
}

// The error stream of the last executed command.
// Null if no command has been executed.
func (r *Container) Stderr() *File {
	q := r.q.Select("stderr")

	return &File{
		q: q,
		c: r.c,
	}
}

// The output stream of the last executed command.
// Null if no command has been executed.
func (r *Container) Stdout() *File {
	q := r.q.Select("stdout")

	return &File{
		q: q,
		c: r.c,
	}
}

// The user to be set for all commands
func (r *Container) User(ctx context.Context) (string, error) {
	q := r.q.Select("user")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// The value of the specified environment variable
func (r *Container) Variable(ctx context.Context, name string) (string, error) {
	q := r.q.Select("variable")
	q = q.Arg("name", name)

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// A list of environment variables passed to commands
func (r *Container) Variables(ctx context.Context) ([]string, error) {
	q := r.q.Select("variables")

	var response []string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// ContainerWithDefaultArgsOpts contains options for Container.WithDefaultArgs
type ContainerWithDefaultArgsOpts struct {
	Args []string
}

// Configures default arguments for future commands
func (r *Container) WithDefaultArgs(opts ...ContainerWithDefaultArgsOpts) *Container {
	q := r.q.Select("withDefaultArgs")
	// `args` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Args) {
			q = q.Arg("args", opts[i].Args)
			break
		}
	}

	return &Container{
		q: q,
		c: r.c,
	}
}

// This container but with a different command entrypoint
func (r *Container) WithEntrypoint(args []string) *Container {
	q := r.q.Select("withEntrypoint")
	q = q.Arg("args", args)

	return &Container{
		q: q,
		c: r.c,
	}
}

// ContainerWithMountedCacheOpts contains options for Container.WithMountedCache
type ContainerWithMountedCacheOpts struct {
	Source DirectoryID
}

// This container plus a cache volume mounted at the given path
func (r *Container) WithMountedCache(cache CacheID, path string, opts ...ContainerWithMountedCacheOpts) *Container {
	q := r.q.Select("withMountedCache")
	q = q.Arg("cache", cache)
	q = q.Arg("path", path)
	// `source` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Source) {
			q = q.Arg("source", opts[i].Source)
			break
		}
	}

	return &Container{
		q: q,
		c: r.c,
	}
}

// This container plus a directory mounted at the given path
func (r *Container) WithMountedDirectory(path string, source DirectoryID) *Container {
	q := r.q.Select("withMountedDirectory")
	q = q.Arg("path", path)
	q = q.Arg("source", source)

	return &Container{
		q: q,
		c: r.c,
	}
}

// This container plus a file mounted at the given path
func (r *Container) WithMountedFile(path string, source FileID) *Container {
	q := r.q.Select("withMountedFile")
	q = q.Arg("path", path)
	q = q.Arg("source", source)

	return &Container{
		q: q,
		c: r.c,
	}
}

// This container plus a secret mounted into a file at the given path
func (r *Container) WithMountedSecret(path string, source SecretID) *Container {
	q := r.q.Select("withMountedSecret")
	q = q.Arg("path", path)
	q = q.Arg("source", source)

	return &Container{
		q: q,
		c: r.c,
	}
}

// This container plus a temporary directory mounted at the given path
func (r *Container) WithMountedTemp(path string) *Container {
	q := r.q.Select("withMountedTemp")
	q = q.Arg("path", path)

	return &Container{
		q: q,
		c: r.c,
	}
}

// This container plus an env variable containing the given secret
func (r *Container) WithSecretVariable(name string, secret SecretID) *Container {
	q := r.q.Select("withSecretVariable")
	q = q.Arg("name", name)
	q = q.Arg("secret", secret)

	return &Container{
		q: q,
		c: r.c,
	}
}

// This container but with a different command user
func (r *Container) WithUser(name string) *Container {
	q := r.q.Select("withUser")
	q = q.Arg("name", name)

	return &Container{
		q: q,
		c: r.c,
	}
}

// This container plus the given environment variable
func (r *Container) WithVariable(name string, value string) *Container {
	q := r.q.Select("withVariable")
	q = q.Arg("name", name)
	q = q.Arg("value", value)

	return &Container{
		q: q,
		c: r.c,
	}
}

// This container but with a different working directory
func (r *Container) WithWorkdir(path string) *Container {
	q := r.q.Select("withWorkdir")
	q = q.Arg("path", path)

	return &Container{
		q: q,
		c: r.c,
	}
}

// This container after unmounting everything at the given path.
func (r *Container) WithoutMount(path string) *Container {
	q := r.q.Select("withoutMount")
	q = q.Arg("path", path)

	return &Container{
		q: q,
		c: r.c,
	}
}

// This container minus the given environment variable
func (r *Container) WithoutVariable(name string) *Container {
	q := r.q.Select("withoutVariable")
	q = q.Arg("name", name)

	return &Container{
		q: q,
		c: r.c,
	}
}

// The working directory for all commands
func (r *Container) Workdir(ctx context.Context) (string, error) {
	q := r.q.Select("workdir")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// Core API
type Core struct {
	q *querybuilder.Selection
	c graphql.Client
}

// Look up a filesystem by its ID
func (r *Core) Filesystem(id FSID) *Filesystem {
	q := r.q.Select("filesystem")
	q = q.Arg("id", id)

	return &Filesystem{
		q: q,
		c: r.c,
	}
}

// CoreGitOpts contains options for Core.Git
type CoreGitOpts struct {
	Ref string
}

func (r *Core) Git(remote string, opts ...CoreGitOpts) *Filesystem {
	q := r.q.Select("git")
	// `ref` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Ref) {
			q = q.Arg("ref", opts[i].Ref)
			break
		}
	}
	q = q.Arg("remote", remote)

	return &Filesystem{
		q: q,
		c: r.c,
	}
}

// Fetch an OCI image
func (r *Core) Image(ref string) *Filesystem {
	q := r.q.Select("image")
	q = q.Arg("ref", ref)

	return &Filesystem{
		q: q,
		c: r.c,
	}
}

// A directory
type Directory struct {
	q *querybuilder.Selection
	c graphql.Client
}

// DirectoryContentsOpts contains options for Directory.Contents
type DirectoryContentsOpts struct {
	Path string
}

// Return a list of files and directories at the given path
func (r *Directory) Contents(ctx context.Context, opts ...DirectoryContentsOpts) ([]string, error) {
	q := r.q.Select("contents")
	// `path` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Path) {
			q = q.Arg("path", opts[i].Path)
			break
		}
	}

	var response []string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// The difference between this directory and an another directory
func (r *Directory) Diff(other DirectoryID) *Directory {
	q := r.q.Select("diff")
	q = q.Arg("other", other)

	return &Directory{
		q: q,
		c: r.c,
	}
}

// Retrieve a directory at the given path
func (r *Directory) Directory(path string) *Directory {
	q := r.q.Select("directory")
	q = q.Arg("path", path)

	return &Directory{
		q: q,
		c: r.c,
	}
}

// Retrieve a file at the given path
func (r *Directory) File(path string) *File {
	q := r.q.Select("file")
	q = q.Arg("path", path)

	return &File{
		q: q,
		c: r.c,
	}
}

// The content-addressed identifier of the directory
func (r *Directory) ID(ctx context.Context) (DirectoryID, error) {
	q := r.q.Select("id")

	var response DirectoryID
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// load a project's metadata
func (r *Directory) LoadProject(configPath string) *Project {
	q := r.q.Select("loadProject")
	q = q.Arg("configPath", configPath)

	return &Project{
		q: q,
		c: r.c,
	}
}

// This directory plus the contents of the given file copied to the given path
func (r *Directory) WithCopiedFile(path string, source FileID) *Directory {
	q := r.q.Select("withCopiedFile")
	q = q.Arg("path", path)
	q = q.Arg("source", source)

	return &Directory{
		q: q,
		c: r.c,
	}
}

// This directory plus a directory written at the given path
func (r *Directory) WithDirectory(directory DirectoryID, path string) *Directory {
	q := r.q.Select("withDirectory")
	q = q.Arg("directory", directory)
	q = q.Arg("path", path)

	return &Directory{
		q: q,
		c: r.c,
	}
}

// DirectoryWithNewFileOpts contains options for Directory.WithNewFile
type DirectoryWithNewFileOpts struct {
	Contents string
}

// This directory plus a new file written at the given path
func (r *Directory) WithNewFile(path string, opts ...DirectoryWithNewFileOpts) *Directory {
	q := r.q.Select("withNewFile")
	// `contents` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Contents) {
			q = q.Arg("contents", opts[i].Contents)
			break
		}
	}
	q = q.Arg("path", path)

	return &Directory{
		q: q,
		c: r.c,
	}
}

// This directory with the directory at the given path removed
func (r *Directory) WithoutDirectory(path string) *Directory {
	q := r.q.Select("withoutDirectory")
	q = q.Arg("path", path)

	return &Directory{
		q: q,
		c: r.c,
	}
}

// This directory with the file at the given path removed
func (r *Directory) WithoutFile(path string) *Directory {
	q := r.q.Select("withoutFile")
	q = q.Arg("path", path)

	return &Directory{
		q: q,
		c: r.c,
	}
}

// Command execution
type Exec struct {
	q *querybuilder.Selection
	c graphql.Client
}

// Exit code of the command
func (r *Exec) ExitCode(ctx context.Context) (int, error) {
	q := r.q.Select("exitCode")

	var response int
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// Modified filesystem
func (r *Exec) FS() *Filesystem {
	q := r.q.Select("fs")

	return &Filesystem{
		q: q,
		c: r.c,
	}
}

// Modified mounted filesystem
func (r *Exec) Mount(path string) *Filesystem {
	q := r.q.Select("mount")
	q = q.Arg("path", path)

	return &Filesystem{
		q: q,
		c: r.c,
	}
}

// ExecStderrOpts contains options for Exec.Stderr
type ExecStderrOpts struct {
	Lines int
}

// stderr of the command
func (r *Exec) Stderr(ctx context.Context, opts ...ExecStderrOpts) (string, error) {
	q := r.q.Select("stderr")
	// `lines` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Lines) {
			q = q.Arg("lines", opts[i].Lines)
			break
		}
	}

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// ExecStdoutOpts contains options for Exec.Stdout
type ExecStdoutOpts struct {
	Lines int
}

// stdout of the command
func (r *Exec) Stdout(ctx context.Context, opts ...ExecStdoutOpts) (string, error) {
	q := r.q.Select("stdout")
	// `lines` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Lines) {
			q = q.Arg("lines", opts[i].Lines)
			break
		}
	}

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// A schema extension provided by a project
type Extension struct {
	q *querybuilder.Selection
	c graphql.Client
}

// path to the extension's code within the project's filesystem
func (r *Extension) Path(ctx context.Context) (string, error) {
	q := r.q.Select("path")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// schema contributed to the project by this extension
func (r *Extension) Schema(ctx context.Context) (string, error) {
	q := r.q.Select("schema")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// sdk used to generate code for and/or execute this extension
func (r *Extension) SDK(ctx context.Context) (string, error) {
	q := r.q.Select("sdk")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// A file
type File struct {
	q *querybuilder.Selection
	c graphql.Client
}

// The contents of the file
func (r *File) Contents(ctx context.Context) (string, error) {
	q := r.q.Select("contents")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// The content-addressed identifier of the file
func (r *File) ID(ctx context.Context) (FileID, error) {
	q := r.q.Select("id")

	var response FileID
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

func (r *File) Secret() *Secret {
	q := r.q.Select("secret")

	return &Secret{
		q: q,
		c: r.c,
	}
}

// The size of the file, in bytes
func (r *File) Size(ctx context.Context) (int, error) {
	q := r.q.Select("size")

	var response int
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// A reference to a filesystem tree.
//
// For example:
//   - The root filesystem of a container
//   - A source code repository
//   - A directory containing binary artifacts
//
// Rule of thumb: if it fits in a tar archive, it fits in a Filesystem.
type Filesystem struct {
	q *querybuilder.Selection
	c graphql.Client
}

// FilesystemCopyOpts contains options for Filesystem.Copy
type FilesystemCopyOpts struct {
	DestPath string

	Exclude []string

	Include []string

	SrcPath string
}

// copy from a filesystem
func (r *Filesystem) Copy(from FSID, opts ...FilesystemCopyOpts) *Filesystem {
	q := r.q.Select("copy")
	// `destPath` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].DestPath) {
			q = q.Arg("destPath", opts[i].DestPath)
			break
		}
	}
	// `exclude` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Exclude) {
			q = q.Arg("exclude", opts[i].Exclude)
			break
		}
	}
	q = q.Arg("from", from)
	// `include` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Include) {
			q = q.Arg("include", opts[i].Include)
			break
		}
	}
	// `srcPath` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].SrcPath) {
			q = q.Arg("srcPath", opts[i].SrcPath)
			break
		}
	}

	return &Filesystem{
		q: q,
		c: r.c,
	}
}

// FilesystemDockerbuildOpts contains options for Filesystem.Dockerbuild
type FilesystemDockerbuildOpts struct {
	Dockerfile string
}

// docker build using this filesystem as context
func (r *Filesystem) Dockerbuild(opts ...FilesystemDockerbuildOpts) *Filesystem {
	q := r.q.Select("dockerbuild")
	// `dockerfile` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Dockerfile) {
			q = q.Arg("dockerfile", opts[i].Dockerfile)
			break
		}
	}

	return &Filesystem{
		q: q,
		c: r.c,
	}
}

// execute a command inside this filesystem
func (r *Filesystem) Exec(input ExecInput) *Exec {
	q := r.q.Select("exec")
	q = q.Arg("input", input)

	return &Exec{
		q: q,
		c: r.c,
	}
}

// FilesystemFileOpts contains options for Filesystem.File
type FilesystemFileOpts struct {
	Lines int
}

// read a file at path
func (r *Filesystem) File(ctx context.Context, path string, opts ...FilesystemFileOpts) (string, error) {
	q := r.q.Select("file")
	// `lines` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Lines) {
			q = q.Arg("lines", opts[i].Lines)
			break
		}
	}
	q = q.Arg("path", path)

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

func (r *Filesystem) ID(ctx context.Context) (FSID, error) {
	q := r.q.Select("id")

	var response FSID
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// push a filesystem as an image to a registry
func (r *Filesystem) PushImage(ctx context.Context, ref string) (bool, error) {
	q := r.q.Select("pushImage")
	q = q.Arg("ref", ref)

	var response bool
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// FilesystemWriteFileOpts contains options for Filesystem.WriteFile
type FilesystemWriteFileOpts struct {
	Permissions string
}

// write a file at path
func (r *Filesystem) WriteFile(contents string, path string, opts ...FilesystemWriteFileOpts) *Filesystem {
	q := r.q.Select("writeFile")
	q = q.Arg("contents", contents)
	q = q.Arg("path", path)
	// `permissions` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Permissions) {
			q = q.Arg("permissions", opts[i].Permissions)
			break
		}
	}

	return &Filesystem{
		q: q,
		c: r.c,
	}
}

// A git ref (tag or branch)
type GitRef struct {
	q *querybuilder.Selection
	c graphql.Client
}

// The digest of the current value of this ref
func (r *GitRef) Digest(ctx context.Context) (string, error) {
	q := r.q.Select("digest")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// The filesystem tree at this ref
func (r *GitRef) Tree() *Directory {
	q := r.q.Select("tree")

	return &Directory{
		q: q,
		c: r.c,
	}
}

// A git repository
type GitRepository struct {
	q *querybuilder.Selection
	c graphql.Client
}

// Details on one branch
func (r *GitRepository) Branch(name string) *GitRef {
	q := r.q.Select("branch")
	q = q.Arg("name", name)

	return &GitRef{
		q: q,
		c: r.c,
	}
}

// List of branches on the repository
func (r *GitRepository) Branches(ctx context.Context) ([]string, error) {
	q := r.q.Select("branches")

	var response []string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// Details on one tag
func (r *GitRepository) Tag(name string) *GitRef {
	q := r.q.Select("tag")
	q = q.Arg("name", name)

	return &GitRef{
		q: q,
		c: r.c,
	}
}

// List of tags on the repository
func (r *GitRepository) Tags(ctx context.Context) ([]string, error) {
	q := r.q.Select("tags")

	var response []string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// Information about the host execution environment
type Host struct {
	q *querybuilder.Selection
	c graphql.Client
}

// Access a directory on the host
func (r *Host) Directory(id HostDirectoryID) *HostDirectory {
	q := r.q.Select("directory")
	q = q.Arg("id", id)

	return &HostDirectory{
		q: q,
		c: r.c,
	}
}

// Lookup the value of an environment variable. Null if the variable is not available.
func (r *Host) Variable(name string) *HostVariable {
	q := r.q.Select("variable")
	q = q.Arg("name", name)

	return &HostVariable{
		q: q,
		c: r.c,
	}
}

// The current working directory on the host
func (r *Host) Workdir() *HostDirectory {
	q := r.q.Select("workdir")

	return &HostDirectory{
		q: q,
		c: r.c,
	}
}

// A directory on the host
type HostDirectory struct {
	q *querybuilder.Selection
	c graphql.Client
}

// Read the contents of the directory
func (r *HostDirectory) Read() *Directory {
	q := r.q.Select("read")

	return &Directory{
		q: q,
		c: r.c,
	}
}

// HostDirectoryWriteOpts contains options for HostDirectory.Write
type HostDirectoryWriteOpts struct {
	Path string
}

// Write the contents of another directory to the directory
func (r *HostDirectory) Write(ctx context.Context, contents DirectoryID, opts ...HostDirectoryWriteOpts) (bool, error) {
	q := r.q.Select("write")
	q = q.Arg("contents", contents)
	// `path` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].Path) {
			q = q.Arg("path", opts[i].Path)
			break
		}
	}

	var response bool
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// An environment variable on the host environment
type HostVariable struct {
	q *querybuilder.Selection
	c graphql.Client
}

// A secret referencing the value of this variable
func (r *HostVariable) Secret() *Secret {
	q := r.q.Select("secret")

	return &Secret{
		q: q,
		c: r.c,
	}
}

// The value of this variable
func (r *HostVariable) Value(ctx context.Context) (string, error) {
	q := r.q.Select("value")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// A set of scripts and/or extensions
type Project struct {
	q *querybuilder.Selection
	c graphql.Client
}

// other projects with schema this project depends on
func (r *Project) Dependencies(ctx context.Context) ([]Project, error) {
	q := r.q.Select("dependencies")

	var response []Project
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// extensions in this project
func (r *Project) Extensions(ctx context.Context) ([]Extension, error) {
	q := r.q.Select("extensions")

	var response []Extension
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// Code files generated by the SDKs in the project
func (r *Project) GeneratedCode() *Directory {
	q := r.q.Select("generatedCode")

	return &Directory{
		q: q,
		c: r.c,
	}
}

// install the project's schema
func (r *Project) Install(ctx context.Context) (bool, error) {
	q := r.q.Select("install")

	var response bool
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// name of the project
func (r *Project) Name(ctx context.Context) (string, error) {
	q := r.q.Select("name")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// schema provided by the project
func (r *Project) Schema(ctx context.Context) (string, error) {
	q := r.q.Select("schema")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// scripts in this project
func (r *Project) Scripts(ctx context.Context) ([]Script, error) {
	q := r.q.Select("scripts")

	var response []Script
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

type Query struct {
	q *querybuilder.Selection
	c graphql.Client
}

// Construct a cache volume from its ID
func (r *Query) Cache(id CacheID) *CacheVolume {
	q := r.q.Select("cache")
	q = q.Arg("id", id)

	return &CacheVolume{
		q: q,
		c: r.c,
	}
}

// Create a new cache volume identified by an arbitrary set of tokens
func (r *Query) CacheFromTokens(tokens []string) *CacheVolume {
	q := r.q.Select("cacheFromTokens")
	q = q.Arg("tokens", tokens)

	return &CacheVolume{
		q: q,
		c: r.c,
	}
}

// ContainerOpts contains options for Query.Container
type ContainerOpts struct {
	ID ContainerID
}

// Load a container from ID.
// Null ID returns an empty container (scratch).
func (r *Query) Container(opts ...ContainerOpts) *Container {
	q := r.q.Select("container")
	// `id` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].ID) {
			q = q.Arg("id", opts[i].ID)
			break
		}
	}

	return &Container{
		q: q,
		c: r.c,
	}
}

// Core API
func (r *Query) Core() *Core {
	q := r.q.Select("core")

	return &Core{
		q: q,
		c: r.c,
	}
}

// DirectoryOpts contains options for Query.Directory
type DirectoryOpts struct {
	ID DirectoryID
}

// Load a directory by ID. No argument produces an empty directory.
func (r *Query) Directory(opts ...DirectoryOpts) *Directory {
	q := r.q.Select("directory")
	// `id` optional argument
	for i := len(opts) - 1; i >= 0; i-- {
		if !querybuilder.IsZeroValue(opts[i].ID) {
			q = q.Arg("id", opts[i].ID)
			break
		}
	}

	return &Directory{
		q: q,
		c: r.c,
	}
}

// Load a file by ID
func (r *Query) File(id FileID) *File {
	q := r.q.Select("file")
	q = q.Arg("id", id)

	return &File{
		q: q,
		c: r.c,
	}
}

// Query a git repository
func (r *Query) Git(url string) *GitRepository {
	q := r.q.Select("git")
	q = q.Arg("url", url)

	return &GitRepository{
		q: q,
		c: r.c,
	}
}

// Query the host environment
func (r *Query) Host() *Host {
	q := r.q.Select("host")

	return &Host{
		q: q,
		c: r.c,
	}
}

// An http remote
func (r *Query) HTTP(url string) *File {
	q := r.q.Select("http")
	q = q.Arg("url", url)

	return &File{
		q: q,
		c: r.c,
	}
}

// Look up a project by name
func (r *Query) Project(name string) *Project {
	q := r.q.Select("project")
	q = q.Arg("name", name)

	return &Project{
		q: q,
		c: r.c,
	}
}

// Load a secret from its ID
func (r *Query) Secret(id SecretID) *Secret {
	q := r.q.Select("secret")
	q = q.Arg("id", id)

	return &Secret{
		q: q,
		c: r.c,
	}
}

// An executable script that uses the project's dependencies and/or extensions
type Script struct {
	q *querybuilder.Selection
	c graphql.Client
}

// path to the script's code within the project's filesystem
func (r *Script) Path(ctx context.Context) (string, error) {
	q := r.q.Select("path")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// sdk used to generate code for and/or execute this script
func (r *Script) SDK(ctx context.Context) (string, error) {
	q := r.q.Select("sdk")

	var response string
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}

// A reference to a secret value, which can be handled more safely than the value itself
type Secret struct {
	q *querybuilder.Selection
	c graphql.Client
}

// The identifier for this secret
func (r *Secret) ID(ctx context.Context) (SecretID, error) {
	q := r.q.Select("id")

	var response SecretID
	q = q.Bind(&response)
	return response, q.Execute(ctx, r.c)
}
