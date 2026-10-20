package main

import (
	"errors"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
// The returned code is the exit code of the child process.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec // the command is provided by the user
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = buildEnv(env)

	err := command.Run()
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	return 1
}

func buildEnv(env Environment) []string {
	values := make(map[string]string)
	for _, pair := range os.Environ() {
		key, value, ok := strings.Cut(pair, "=")
		if !ok {
			continue
		}
		values[key] = value
	}

	for name, value := range env {
		if strings.ContainsRune(name, '=') {
			continue
		}

		if value.NeedRemove {
			delete(values, name)
			continue
		}

		values[name] = value.Value
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, key+"="+values[key])
	}

	return result
}
