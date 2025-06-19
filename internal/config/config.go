// config.go
// Environment and config loading for EMSG Daemon
package main

import (
	"bufio"
	"os"
	"strings"
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

// LoadConfigFromFile loads config from a .env or config file
func LoadConfigFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	cfg := &Config{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "EMSG_DATABASE_URL":
			cfg.DatabaseURL = val
		case "EMSG_DOMAIN":
			cfg.Domain = val
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cfg, nil
}
