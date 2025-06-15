package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	dirInfo, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envVars := make(Environment)

	s, _ := os.OpenFile("a", 0, 0)

	var a io.ReadWriter = s

	var any interface{} = a

	any.(*os.File).Name()

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
