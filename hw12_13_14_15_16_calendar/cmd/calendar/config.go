package main

import (
	"fmt"
	"github.com/davilov/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/sql"
	"gopkg.in/yaml.v3"
	"os"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger      logger.LogConf          `yaml:"logger"`
	Database    sqlstorage.DatabaseConf `yaml:"postgres"`
	StorageType string                  `yaml:"storage_type"`
	// TODO
}

func NewConfig() Config {
	_, err := os.Stat(configFile)
	if err != nil {
		fmt.Println("Error: File does not exist")
		os.Exit(1)
	}
	file, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error: File could not be read")
		os.Exit(1)
	}

	cnf := Config{}
	err = yaml.Unmarshal(file, &cnf)
	if err != nil {
		fmt.Println("Error: File could not be parsed")
		os.Exit(1)
	}

	return cnf
}

// TODO
