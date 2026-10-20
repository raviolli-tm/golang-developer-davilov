package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("whole file", func(t *testing.T) {
		target := filepath.Join(t.TempDir(), "out.txt")
		require.NoError(t, Copy("testdata/input.txt", target, 0, 0))

		expected, err := os.ReadFile("testdata/out_offset0_limit0.txt")
		require.NoError(t, err)
		actual, err := os.ReadFile(target)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("offset and limit", func(t *testing.T) {
		target := filepath.Join(t.TempDir(), "out.txt")
		require.NoError(t, Copy("testdata/input.txt", target, 100, 1000))

		expected, err := os.ReadFile("testdata/out_offset100_limit1000.txt")
		require.NoError(t, err)
		actual, err := os.ReadFile(target)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("limit larger than file", func(t *testing.T) {
		target := filepath.Join(t.TempDir(), "out.txt")
		require.NoError(t, Copy("testdata/input.txt", target, 0, 1_000_000))

		expected, err := os.ReadFile("testdata/out_offset0_limit0.txt")
		require.NoError(t, err)
		actual, err := os.ReadFile(target)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("offset equal to file size", func(t *testing.T) {
		info, err := os.Stat("testdata/input.txt")
		require.NoError(t, err)

		target := filepath.Join(t.TempDir(), "out.txt")
		require.NoError(t, Copy("testdata/input.txt", target, info.Size(), 0))

		actual, err := os.ReadFile(target)
		require.NoError(t, err)
		require.Empty(t, actual)
	})
}

func TestCopyEdgeErrors(t *testing.T) {
	t.Run("file paths are required", func(t *testing.T) {
		require.ErrorIs(t, Copy("", "out.txt", 0, 0), ErrFilePathIsRequired)
		require.ErrorIs(t, Copy("testdata/input.txt", "", 0, 0), ErrFilePathIsRequired)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		target := filepath.Join(t.TempDir(), "out.txt")
		require.ErrorIs(t, Copy("testdata/input.txt", target, 1_000_000, 0), ErrOffsetExceedsFileSize)
	})

	t.Run("negative offset and limit", func(t *testing.T) {
		require.ErrorIs(t, Copy("testdata/input.txt", filepath.Join(t.TempDir(), "out.txt"), -1, 0), ErrInvalidOffset)
		require.ErrorIs(t, Copy("testdata/input.txt", filepath.Join(t.TempDir(), "out.txt"), 0, -1), ErrInvalidLimit)
	})

	t.Run("directory is not supported", func(t *testing.T) {
		require.ErrorIs(t, Copy("testdata", filepath.Join(t.TempDir(), "out.txt"), 0, 0), ErrUnsupportedFile)
	})

	t.Run("missing source", func(t *testing.T) {
		require.Error(t, Copy("testdata/missing.txt", filepath.Join(t.TempDir(), "out.txt"), 0, 0))
	})

	t.Run("invalid target", func(t *testing.T) {
		require.Error(t, Copy("testdata/input.txt", filepath.Join(t.TempDir(), "missing", "out.txt"), 0, 0))
	})
}
