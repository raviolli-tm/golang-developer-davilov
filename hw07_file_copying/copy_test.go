package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyErrors(t *testing.T) {
	t.Run("Check errors", func(t *testing.T) {
		err := Copy("", "a.txt", 10, 10)
		if err != nil {
			assert.ErrorIs(t, err, ErrFilePathIsRequired)
		}

		err = Copy("a.txt", "", 10, 10)
		if err != nil {
			assert.ErrorIs(t, err, ErrFilePathIsRequired)
		}

		err = Copy("testdata/input.txt", "out.txt", 20000, 10)
		if err != nil {
			assert.ErrorIs(t, err, ErrOffsetExceedsFileSize)
		}
	})
}
