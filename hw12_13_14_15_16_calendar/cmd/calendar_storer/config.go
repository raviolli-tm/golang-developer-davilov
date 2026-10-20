package main

import (
	"fmt"

	"github.com/davilov/hw12_13_14_15_calendar/internal/broker"
	"github.com/davilov/hw12_13_14_15_calendar/internal/config"
	sqlstorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/sql"
)

type Config struct {
	Database sqlstorage.DatabaseConf `yaml:"postgres"`
	Broker   broker.KafkaConf        `yaml:"broker"`
	Metrics  MetricsConf             `yaml:"metrics"`
}

type MetricsConf struct {
	Port int `yaml:"port"`
}

func NewConfig() (Config, error) {
	cnf := Config{}

	if err := config.Load(configFile, &cnf); err != nil {
		return Config{}, err
	}

	return cnf, nil
}

func (c Config) MetricsAddr() string {
	return fmt.Sprintf(":%d", c.Metrics.Port)
}
