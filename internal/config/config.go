// config.go
// Environment and config loading for EMSG Daemon
package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL    string
	Domain         string
	Port           string
	LogLevel       string
	MaxConnections int
	WWWDir         string
	LogFile        string
}

// binaryDir returns the directory containing the running binary.
func binaryDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

func LoadConfig() (*Config, error) {
	dir := binaryDir()
	// Render (and many PaaS platforms) inject PORT; fall back to EMSG_PORT then 8765.
	port := getEnvWithDefault("PORT", getEnvWithDefault("EMSG_PORT", "8765"))
	cfg := &Config{
		DatabaseURL:    getEnvWithDefault("EMSG_DATABASE_URL", ""),
		Domain:         getEnvWithDefault("EMSG_DOMAIN", ""),
		Port:           port,
		LogLevel:       getEnvWithDefault("EMSG_LOG_LEVEL", "info"),
		MaxConnections: getEnvIntWithDefault("EMSG_MAX_CONNECTIONS", 100),
		WWWDir:         getEnvWithDefault("EMSG_WWW_DIR", filepath.Join(dir, "www")),
		LogFile:        getEnvWithDefault("EMSG_LOG_FILE", filepath.Join(dir, "emsg.log")),
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
		if n, err := strconv.Atoi(value); err == nil {
			return n
		}
	}
	return defaultValue
}

// LoadConfigFromFile loads config from a .env or config file.
// Returns an error if the file cannot be opened (e.g. absent), so the caller
// can log a warning and fall back to defaults.
func LoadConfigFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dir := binaryDir()
	cfg := &Config{
		// Pre-populate with binary-relative defaults so partial config files
		// still produce a fully-initialised struct.
		Port:           "8765",
		LogLevel:       "info",
		MaxConnections: 100,
		WWWDir:         filepath.Join(dir, "www"),
		LogFile:        filepath.Join(dir, "emsg.log"),
	}

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
			if n, err := strconv.Atoi(val); err == nil {
				cfg.MaxConnections = n
			}
		case "EMSG_WWW_DIR":
			cfg.WWWDir = val
		case "EMSG_LOG_FILE":
			cfg.LogFile = val
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cfg, nil
}
