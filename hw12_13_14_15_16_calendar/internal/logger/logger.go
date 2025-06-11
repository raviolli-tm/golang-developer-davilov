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
	intLevel := 0
	switch loggerConf.Level {
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

	destinations := make([]io.Writer, 0)
	if loggerConf.Stdout.Enabled {
		destinations = append(destinations, os.Stdout)
	}

	loggerUniqueDateName := time.Now().Format("2006_01_02_150405")

	if loggerConf.File.Enabled {
		for _, destination := range loggerConf.File.Destinations {
			filePath := destination + "log_" + loggerUniqueDateName + ".log"
			destinationFile, fileErr := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if fileErr != nil {
				fmt.Printf("ERROR | Error opening log file: %s \n", filePath)
			}
			destinations = append(destinations, destinationFile)
		}
	}

	log := &Logger{loggerConf.Level, intLevel, destinations}

	return log
}

const dateFormat string = "01-02-2006 15:04:05.000"
const format string = "%s [%s] %s \n"

func (l *Logger) Warn(msg string) {
	currentTime := time.Now().Format(dateFormat)
	for _, destination := range l.destinations {
		err(msg, destination, currentTime)
	}
}

func (l *Logger) Info(msg string) {
	currentTime := time.Now().Format(dateFormat)
	for _, destination := range l.destinations {
		info(msg, destination, currentTime)
	}
}

func (l *Logger) Debug(msg string) {
	currentTime := time.Now().Format(dateFormat)
	for _, destination := range l.destinations {
		debug(msg, destination, currentTime)
	}
}

func (l *Logger) Error(msg string) {
	currentTime := time.Now().Format(dateFormat)
	for _, destination := range l.destinations {
		err(msg, destination, currentTime)
	}
}

func warn(msg string, destination io.Writer, time string) {
	logMsg(msg, destination, time, "WARN")
}

func info(msg string, destination io.Writer, time string) {
	logMsg(msg, destination, time, "INFO")
}

func debug(msg string, destination io.Writer, time string) {
	logMsg(msg, destination, time, "DEBUG")
}

func err(msg string, destination io.Writer, time string) {
	logMsg(msg, destination, time, "ERROR")
}
func logMsg(msg string, destination io.Writer, time string, logLevel string) {
	_, err := destination.Write([]byte(fmt.Sprintf(format, time, logLevel, msg)))
	if err != nil {
		return
	}
}
