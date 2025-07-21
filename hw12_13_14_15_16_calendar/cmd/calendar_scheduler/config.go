package main

import (
	"fmt"
	"github.com/davilov/hw12_13_14_15_calendar/internal/broker"
	sqlstorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/sql"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Database    sqlstorage.DatabaseConf `yaml:"postgres"`
	Broker      broker.KafkaConf        `yaml:"broker"`
	RetryPolicy broker.RetryPolicyConf  `yaml:"retry"`
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
