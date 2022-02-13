package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Host  string
	Port  int
	Debug bool
}

const (
	configName = "config"
	configType = "yaml"
)

func NewConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("empty path")
	}
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failure to read config directory %s: %s", configPath, err)
	}
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failture to unmarshal config into struct: %s", err)
	}
	return cfg, nil
}
