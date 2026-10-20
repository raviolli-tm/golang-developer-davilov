package main

import (
	"github.com/davilov/hw12_13_14_15_calendar/internal/broker"
	"github.com/davilov/hw12_13_14_15_calendar/internal/config"
	sqlstorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/sql"
)

type Config struct {
	Database             sqlstorage.DatabaseConf `yaml:"postgres"`
	Broker               broker.KafkaConf        `yaml:"broker"`
	RetryPolicy          broker.RetryPolicyConf  `yaml:"retry"`
	CleanupIntervalHours int                     `yaml:"cleanup_interval_hours"`
}

func NewConfig() (Config, error) {
	cnf := Config{}

	if err := config.Load(configFile, &cnf); err != nil {
		return Config{}, err
	}

	return cnf, nil
}
