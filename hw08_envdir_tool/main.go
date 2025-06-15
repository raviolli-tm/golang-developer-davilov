package main

import (
	"fmt"
	"os"
)

func main() {
	envPath := os.Args[1]

	cmd := make([]string, 0)

	cmd = append(cmd, os.Args[2:]...)

	env, err := ReadDir(envPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	RunCmd(cmd, env)
}
