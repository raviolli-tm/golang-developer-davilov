package main

import (
	"bytes"
	"io"
	"os"
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
		file, err := os.OpenFile(dir+"/"+fileInfo.Name(), os.O_RDONLY, 0666)
		if err != nil {
			return nil, err
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
			fileEnv.Value = strings.TrimRightFunc(fileEnv.Value, func(r rune) bool { return unicode.IsSpace(r) })
		}

		envVars[fileInfo.Name()] = fileEnv
		_ = file.Close()
	}

	return envVars, nil
}
