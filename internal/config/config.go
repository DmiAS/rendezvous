package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host  string
	Port  string
	Debug bool
}

const (
	portEnv  = "PORT"
	debugEnv = "DEBUG"
)

func NewConfig() *Config {
	// load env from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("failure to load env")
	}

	cfg := &Config{
		Host:  "localhost",
		Port:  os.Getenv(portEnv),
		Debug: os.Getenv(debugEnv) != "",
	}
	return cfg
}
