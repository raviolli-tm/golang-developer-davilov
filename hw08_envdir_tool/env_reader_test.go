package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func runCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func permissionsReset(filename string) error {
	switch runtime.GOOS {
	case "linux", "darwin", "freebsd", "openbsd":
		return os.Chmod(filename, 0o644)
	case "windows":
		return runCommand("icacls", filename, "/reset")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func permissionsNoAccess(filename string) error {
	switch runtime.GOOS {
	case "linux", "darwin", "freebsd", "openbsd":
		return os.Chmod(filename, 0o000)
	case "windows":
		_ = runCommand("icacls", filename, "/reset")
		return runCommand("icacls", filename, "/deny", "Everyone:F")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

const (
	testDataDir  = "./testdata/env_test" // Используйте "testdata" для тестовых данных
	testFileName = "test_file_for_test"  // snake_case для имен файлов
)

func TestReadDir(t *testing.T) {
	_ = os.Mkdir(testDataDir, 0o755)
	_, err := os.Create(testDataDir + "/" + testFileName)
	if err != nil {
		fmt.Println("Fatal error: 3333 ", err)
	}

	t.Run("ReadDir Errors", func(t *testing.T) {
		t.Run("bad dir path", func(t *testing.T) {
			_, err := ReadDir("./wewswe")
			if err != nil {
				t.Errorf("ReadDir Errors: %v", err)
			}
		})

		t.Run("file all rights deny", func(t *testing.T) {
			err := permissionsNoAccess(testDataDir)
			if err != nil {
				fmt.Println(err)
			}
			_, err = ReadDir(testDataDir)
			if err != nil {
				t.Errorf("ReadDir Errors: %v", err)
			}
		})
		_ = permissionsReset(testDataDir)
		_ = permissionsReset(testDataDir + "/" + testFileName)
		_ = os.Remove(testDataDir + "/" + testFileName)
		_ = os.Remove(testDataDir)
	})
}
