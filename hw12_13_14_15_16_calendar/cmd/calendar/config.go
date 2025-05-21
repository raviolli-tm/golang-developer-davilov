package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf   `yaml:"logger"`
	Database DatabaseConf `yaml:"postgres"`
	// TODO
}

type LoggerConf struct {
	Level string `yaml:"level"`
	// TODO
}

type DatabaseConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
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
