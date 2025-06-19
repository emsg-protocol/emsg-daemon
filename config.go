// config.go
// Environment and config loading for EMSG Daemon
package main

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Domain      string
}

func LoadConfig() (*Config, error) {
	// TODO: Load config from environment or file
	return &Config{
		DatabaseURL: os.Getenv("EMSG_DATABASE_URL"),
		Domain:      os.Getenv("EMSG_DOMAIN"),
	}, nil
}
