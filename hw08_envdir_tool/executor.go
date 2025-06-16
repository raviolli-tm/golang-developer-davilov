package main

import (
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 || len(cmd) == 1 {
		return 0
	}

	for varName, varValue := range env {
		if varValue.NeedRemove {
			_ = os.Unsetenv(varName)
		} else {
			_ = os.Setenv(varName, varValue.Value)
		}
	}
	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = os.Environ()

	err := command.Run()
	if err != nil {
		return 0
	}

	return 1
}
