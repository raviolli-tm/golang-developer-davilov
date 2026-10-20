package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %q: %w", dir, err)
	}

	env := make(Environment)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.ContainsRune(name, '=') {
			continue
		}

		value, needRemove, err := readEnvFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}

		env[name] = EnvValue{Value: value, NeedRemove: needRemove}
	}

	return env, nil
}

func readEnvFile(path string) (string, bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", false, fmt.Errorf("open env file %q: %w", path, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return "", false, fmt.Errorf("stat env file %q: %w", path, err)
	}
	if info.Size() == 0 {
		return "", true, nil
	}

	line, err := bufio.NewReader(file).ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", false, fmt.Errorf("read env file %q: %w", path, err)
	}

	line = bytes.TrimSuffix(line, []byte{'\n'})
	line = bytes.ReplaceAll(line, []byte{0}, []byte{'\n'})

	return strings.TrimRight(string(line), " \t"), false, nil
}
