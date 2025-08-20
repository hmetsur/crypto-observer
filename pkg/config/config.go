// pkg/config/config.go
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Addr string `yaml:"addr"` // ":8080"
	} `yaml:"server"`

	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`

	Collector struct {
		DefaultPeriodSeconds int `yaml:"default_period_seconds"` // N секунд по умолчанию
	} `yaml:"collector"`

	Coingecko struct {
		BaseURL    string `yaml:"base_url"`  // https://api.coingecko.com/api/v3
		TimeoutSec int    `yaml:"timeout_s"` // 5
	} `yaml:"coingecko"`

	Log struct {
		Level string `yaml:"level"` // info|debug|warn|error
	} `yaml:"log"`
}

var cfg Config

func MustLoad() *Config {
	path := "configs/config.yaml"
	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		panic(err)
	}
	return &cfg
}

// экспортируем геттеры, чтобы другие пакеты не тащили yaml напрямую
func C() *Config { return &cfg }
