package main

import (
	"fmt"

	"github.com/davilov/hw12_13_14_15_calendar/internal/config"
	"github.com/davilov/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/sql"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger      logger.LogConf          `yaml:"logger"`
	Database    sqlstorage.DatabaseConf `yaml:"postgres"`
	StorageType string                  `yaml:"storage_type"`
	HTTP        HTTPConf                `yaml:"http"`
	Metrics     MetricsConf             `yaml:"metrics"`
}

type HTTPConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
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

func (c Config) HTTPAddr() string {
	return fmt.Sprintf("%s:%d", c.HTTP.Host, c.HTTP.Port)
}

func (c Config) MetricsAddr() string {
	return fmt.Sprintf(":%d", c.Metrics.Port)
}
