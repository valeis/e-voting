package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Database      Postgres `yaml:"postgres"`
	RedisDatabase Redis    `yaml:"redisDatabase"`
}

type Postgres struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	RedisDB  int    `yaml:"db"`
}

var Cfg Config

func init() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Errorf("Failed to get working directory: %w", err)
	}
	configPath := filepath.Join(wd, "config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Errorf("Failed to read configuration file: %w", err)
	}
	if err := yaml.Unmarshal(data, &Cfg); err != nil {
		fmt.Errorf("Failed to unmarshal YAML data: %w", err)
	}
}
