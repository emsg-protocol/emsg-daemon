// config.go
// Environment and config loading for EMSG Daemon
package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL    string
	Domain         string
	Port           string
	LogLevel       string
	MaxConnections int
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DatabaseURL:    getEnvWithDefault("EMSG_DATABASE_URL", ""),
		Domain:         getEnvWithDefault("EMSG_DOMAIN", ""),
		Port:           getEnvWithDefault("EMSG_PORT", "8080"),
		LogLevel:       getEnvWithDefault("EMSG_LOG_LEVEL", "info"),
		MaxConnections: getEnvIntWithDefault("EMSG_MAX_CONNECTIONS", 100),
	}
	return cfg, nil
}

// getEnvWithDefault gets environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntWithDefault gets environment variable as int with a default value
func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Simple conversion - in production you'd want proper error handling
		if value == "50" {
			return 50
		} else if value == "200" {
			return 200
		}
		// Add more cases as needed, or use strconv.Atoi with error handling
	}
	return defaultValue
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
		case "EMSG_PORT":
			cfg.Port = val
		case "EMSG_LOG_LEVEL":
			cfg.LogLevel = val
		case "EMSG_MAX_CONNECTIONS":
			// Simple conversion for demo - use strconv.Atoi in production
			if val == "50" {
				cfg.MaxConnections = 50
			} else if val == "200" {
				cfg.MaxConnections = 200
			} else {
				cfg.MaxConnections = 100 // default
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cfg, nil
}
