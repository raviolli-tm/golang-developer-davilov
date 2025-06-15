package main

import (
	"fmt"
	"os"
)

func main() {
	// Place your code here.
	envPath := os.Args[1]

	cmd := make([]string, 0)

	for _, env := range os.Args[2:] {
		cmd = append(cmd, env)
	}

	env, err := ReadDir(envPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	RunCmd(cmd, env)

}
