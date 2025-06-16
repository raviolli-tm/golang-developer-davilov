package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	dirInfo, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envVars := make(Environment)

	for _, fileInfo := range dirInfo {
		file, err := os.Open(
			filepath.Join(dir, fileInfo.Name()))
		if err != nil {
			return nil, fmt.Errorf("error opening file %s : %w", fileInfo.Name(), err)
		}

		fileEnv := EnvValue{"", false}
		all, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		if len(all) == 0 {
			fileEnv.NeedRemove = true
		} else {
			newLineIdx := bytes.IndexRune(all, '\n')
			if newLineIdx != -1 {
				all = all[:newLineIdx]
			}
			all = bytes.ReplaceAll(all, []byte{0}, []byte{10})
			fileEnv.Value = string(all)
			fileEnv.Value = strings.TrimRightFunc(fileEnv.Value, unicode.IsSpace)
		}

		envVars[fileInfo.Name()] = fileEnv
		_ = file.Close()
	}

	return envVars, nil
}
