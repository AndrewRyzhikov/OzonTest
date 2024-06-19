package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DataBaseConfig   DataBaseConfig   `yaml:"data_base" env-required:"true"`
	HttpServerConfig HttpServerConfig `yaml:"http_server" env-required:"true"`
	LogConfig        LogConfig        `yaml:"log" env-required:"true"`
}

type HttpServerConfig struct {
	Port string `yaml:"port" env-required:"true"`
}

type DataBaseConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DbName   string `yaml:"dbname" env-required:"true"`
	Driver   string `yaml:"driver" env-required:"true"`
}

type LogConfig struct {
	Level      string           `yaml:"level" env-default:"INFO"`
	Path       string           `yaml:"path" env-required:"true"`
	Lumberjack LumberjackConfig `yaml:"lumberjack" env-required:"true"`
}

type LumberjackConfig struct {
	MaxSize    uint64 `yaml:"max_size"`
	MaxAge     uint64 `yaml:"max_age"`
	MaxBackups uint64 `yaml:"max_backups"`
	LocalTime  bool   `yaml:"local_time"`
	Compress   bool   `yaml:"compress"`
}

func Load() (*Config, error) {
	path := fetchConfig()
	if path == "" {
		return nil, errors.New("config file not exist")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config path does not exist: %s", path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}

func fetchConfig() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config")
	flag.Parse()

	return res
}
