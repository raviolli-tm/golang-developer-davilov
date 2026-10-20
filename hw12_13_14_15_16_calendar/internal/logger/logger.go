package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Logger struct {
	level        string
	intLevel     int
	destinations []io.Writer
}

type LogConf struct {
	Level  string `yaml:"level"`
	Stdout struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"stdout"`
	File struct {
		Enabled      bool     `yaml:"enabled"`
		Destinations []string `yaml:"destination"`
	} `yaml:"file"`
	Remote struct {
		Enabled      bool     `yaml:"enabled"`
		Destinations []string `yaml:"destination"`
	} `yaml:"remote"`
}

func New(loggerConf LogConf) *Logger {
	intLevel := levelValue(loggerConf.Level)

	destinations := make([]io.Writer, 0)
	if loggerConf.Stdout.Enabled {
		destinations = append(destinations, os.Stdout)
	}

	loggerUniqueDateName := time.Now().Format("2006_01_02_150405")

	if loggerConf.File.Enabled {
		for _, destination := range loggerConf.File.Destinations {
			filePath := destination + "log_" + loggerUniqueDateName + ".log"
			destinationFile, fileErr := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if fileErr != nil {
				fmt.Printf("ERROR | Error opening log file: %s \n", filePath)
				continue
			}
			destinations = append(destinations, destinationFile)
		}
	}

	return &Logger{loggerConf.Level, intLevel, destinations}
}

func levelValue(level string) int {
	switch level {
	case "ERROR":
		return 3
	case "WARN":
		return 2
	case "INFO":
		return 1
	case "DEBUG":
		return 0
	default:
		return 0
	}
}

const dateFormat string = "01-02-2006 15:04:05.000"
const format string = "%s [%s] %s \n"

func (l *Logger) Warn(msg string) {
	l.log(2, "WARN", msg)
}

func (l *Logger) Info(msg string) {
	l.log(1, "INFO", msg)
}

func (l *Logger) Debug(msg string) {
	l.log(0, "DEBUG", msg)
}

func (l *Logger) Error(msg string) {
	l.log(3, "ERROR", msg)
}

func (l *Logger) log(level int, tag, msg string) {
	if level < l.intLevel {
		return
	}

	currentTime := time.Now().Format(dateFormat)
	for _, destination := range l.destinations {
		logMsg(msg, destination, currentTime, tag)
	}
}

func logMsg(msg string, destination io.Writer, time, logLevel string) {
	_, err := destination.Write([]byte(fmt.Sprintf(format, time, logLevel, msg)))
	if err != nil {
		return
	}
}
