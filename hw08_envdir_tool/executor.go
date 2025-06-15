package main

import (
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {

	for varName, varValue := range env {
		if varValue.NeedRemove {
			_ = os.Unsetenv(varName)
		} else {
			_ = os.Setenv(varName, varValue.Value)
		}

	}

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
