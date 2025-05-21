package logger

import (
	"fmt"
	"time"
)

type Logger struct {
	level    string
	intLevel int
}

func New(level string) *Logger {
	intLevel := 0
	switch level {
	case "ERROR":
		intLevel = 3
	case "WARN":
		intLevel = 2
	case "INFO":
		intLevel = 1
	case "DEBUG":
		intLevel = 0
	default:
		intLevel = 0

	}

	return &Logger{level, intLevel}
}

const format string = "01-02-2006 15:04:05.000"

func (l Logger) Warn(msg string) {
	if l.intLevel <= 2 {
		fmt.Println(time.Now().Format(format) + " | WARN  | " + msg)
	}
}

func (l Logger) Info(msg string) {
	if l.intLevel <= 1 {
		fmt.Println(time.Now().Format(format) + " | INFO  | " + msg)
	}
}

func (l Logger) Debug(msg string) {
	if l.intLevel == 0 {
		fmt.Println(time.Now().Format(format) + " | DEBUG | " + msg)
	}
}

func (l Logger) Error(msg string) {
	fmt.Println(time.Now().Format(format) + " | ERROR | " + msg)
}

// TODO
