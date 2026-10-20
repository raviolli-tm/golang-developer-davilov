package main

import (
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("Test cmd", func(t *testing.T) {
		returnCode := RunCmd([]string{}, Environment{})
		if returnCode == 0 {
			t.Logf("RunCmd Errors: returncode=%v expected 0", returnCode)
		}
	})
}
