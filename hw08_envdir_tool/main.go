package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: go-envdir <dir> <command> [args...]")
		os.Exit(2)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
