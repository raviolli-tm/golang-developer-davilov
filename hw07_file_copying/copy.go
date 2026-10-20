package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	sourceFile, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	size, err := sourceFile.Seek(0, 2)
	defer sourceFile.Close()

	if err != nil {
		return ErrUnsupportedFile
	}
	if size < offset {
		return ErrOffsetExceedsFileSize
	}

	sourceFile.Seek(offset, 0)
	targetFile, _ := os.Create(toPath)
	defer targetFile.Close()
	var n int64

	if size-offset < limit || limit == 0 {
		limit = size - offset
	}

	bar := pb.New64(limit).Start()
	barReader := bar.NewProxyReader(sourceFile)
	n, err = io.CopyN(targetFile, barReader, limit)
	bar.Finish()

	if err != nil {
		return err
	}

	fmt.Printf("%d bytes copied from %s to %s\n", n, from, to)

	return nil
}
