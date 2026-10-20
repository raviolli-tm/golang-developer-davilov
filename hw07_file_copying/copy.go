package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrFilePathIsRequired    = errors.New("filepath is required")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidOffset         = errors.New("offset must not be negative")
	ErrInvalidLimit          = errors.New("limit must not be negative")
)

// Copy copies at most limit bytes from fromPath starting at offset into
// toPath. A zero limit means "copy the rest of the file".
func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrFilePathIsRequired
	}
	if offset < 0 {
		return ErrInvalidOffset
	}
	if limit < 0 {
		return ErrInvalidLimit
	}

	sourceFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer sourceFile.Close()

	info, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("stat source file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	size := info.Size()
	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	remaining := size - offset
	if limit == 0 || limit > remaining {
		limit = remaining
	}

	if _, err := sourceFile.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("seek source file: %w", err)
	}

	targetFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("create target file: %w", err)
	}
	defer targetFile.Close()

	copied, err := io.CopyN(targetFile, sourceFile, limit)
	if err != nil {
		return fmt.Errorf("copy data: %w", err)
	}

	if err := targetFile.Close(); err != nil {
		return fmt.Errorf("close target file: %w", err)
	}

	fmt.Printf("%d bytes copied from %s to %s\n", copied, fromPath, toPath)

	return nil
}
