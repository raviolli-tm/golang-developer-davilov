package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoggerFiltersLevels(t *testing.T) {
	dir := t.TempDir()

	conf := LogConf{Level: "ERROR"}
	conf.File.Enabled = true
	conf.File.Destinations = []string{dir + string(os.PathSeparator)}

	logg := New(conf)
	logg.Debug("debug message")
	logg.Info("info message")
	logg.Warn("warn message")
	logg.Error("error message")

	files, err := filepath.Glob(filepath.Join(dir, "log_*.log"))
	if err != nil {
		t.Fatalf("find log files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected one log file, got %d", len(files))
	}

	content, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}

	message := string(content)
	if strings.Contains(message, "debug message") || strings.Contains(message, "info message") {
		t.Fatalf("messages below the configured level must not be written: %s", message)
	}
	if strings.Contains(message, "warn message") {
		t.Fatalf("warn message must not be written for the ERROR level: %s", message)
	}
	if !strings.Contains(message, "ERROR") || !strings.Contains(message, "error message") {
		t.Fatalf("error message was not written: %s", message)
	}
}

func TestLoggerWritesAllLevels(t *testing.T) {
	dir := t.TempDir()

	conf := LogConf{Level: "DEBUG"}
	conf.Stdout.Enabled = false
	conf.File.Enabled = true
	conf.File.Destinations = []string{dir + string(os.PathSeparator)}

	logg := New(conf)
	logg.Debug("debug message")
	logg.Info("info message")
	logg.Warn("warn message")
	logg.Error("error message")

	files, err := filepath.Glob(filepath.Join(dir, "log_*.log"))
	if err != nil || len(files) != 1 {
		t.Fatalf("expected one log file, got %v (err %v)", files, err)
	}

	content, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}

	for _, expected := range []string{"DEBUG", "INFO", "WARN", "ERROR"} {
		if !strings.Contains(string(content), expected) {
			t.Fatalf("log file does not contain %s level: %s", expected, content)
		}
	}
}
