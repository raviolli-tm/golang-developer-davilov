package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDirExtra(t *testing.T) {
	dir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(dir, "FOO"), []byte("bar\nsecond line"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "EMPTY"), nil, 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "EMPTYLINE"), []byte("\nrest"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SPACES"), []byte("value \t "), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "NUL"), []byte("a\x00b"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "A=B"), []byte("invalid"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "subdir"), 0o755))

	env, err := ReadDir(dir)
	require.NoError(t, err)

	require.Equal(t, EnvValue{Value: "bar"}, env["FOO"])
	require.Equal(t, EnvValue{NeedRemove: true}, env["EMPTY"])
	require.Equal(t, EnvValue{Value: ""}, env["EMPTYLINE"])
	require.Equal(t, EnvValue{Value: "value"}, env["SPACES"])
	require.Equal(t, EnvValue{Value: "a\nb"}, env["NUL"])
	require.NotContains(t, env, "A=B")
	require.NotContains(t, env, "subdir")
}

func TestRunCmdExitCodes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		require.Equal(t, 0, RunCmd([]string{"sh", "-c", "exit 0"}, Environment{}))
	})

	t.Run("child exit code", func(t *testing.T) {
		require.Equal(t, 42, RunCmd([]string{"sh", "-c", "exit 42"}, Environment{}))
	})

	t.Run("empty command", func(t *testing.T) {
		require.Equal(t, 1, RunCmd(nil, Environment{}))
	})

	t.Run("missing command", func(t *testing.T) {
		require.Equal(t, 1, RunCmd([]string{"definitely-not-a-command"}, Environment{}))
	})
}

func TestRunCmdEnvironment(t *testing.T) {
	t.Run("sets variable", func(t *testing.T) {
		env := Environment{"FOO": {Value: "bar"}}
		require.Equal(t, 0, RunCmd([]string{"sh", "-c", `test "$FOO" = bar`}, env))
	})

	t.Run("removes variable", func(t *testing.T) {
		t.Setenv("TO_REMOVE", "value")

		env := Environment{"TO_REMOVE": {NeedRemove: true}}
		require.Equal(t, 0, RunCmd([]string{"sh", "-c", `test -z "$TO_REMOVE"`}, env))
	})
}
