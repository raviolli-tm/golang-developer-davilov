package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
		return os.Chmod(filename, 0644)
	case "windows":
		return runCommand("icacls", filename, "/reset")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func permissionsNoAccess(filename string) error {
	switch runtime.GOOS {
	case "linux", "darwin", "freebsd", "openbsd":
		return os.Chmod(filename, 0000)
	case "windows":
		_ = runCommand("icacls", filename, "/reset")
		return runCommand("icacls", filename, "/deny", "Everyone:F")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func TestReadDir(t *testing.T) {

	var dirPath = "./testData/envTest"
	var fileName = "testFileForTestReason"

	t.Run("ReadDir Errors", func(t *testing.T) {
		_, err := ReadDir("./wewswe")
		if err != nil {
			assert.ErrorIs(t, err, os.ErrNotExist)
		}

		createdFile, _ := os.Create(dirPath + "/" + fileName)
		_ = createdFile.Close()
		err = permissionsNoAccess(dirPath + "/" + fileName)
		if err != nil {
			fmt.Println(err)
		}
		_, err = ReadDir(dirPath + "/" + fileName)
		if err != nil {
			assert.ErrorIs(t, err, os.ErrPermission)
		}
		_ = os.Remove(dirPath + "/" + fileName)

	})
}
