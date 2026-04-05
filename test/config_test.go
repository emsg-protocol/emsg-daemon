// config_test.go
// Tests for config loading in EMSG Daemon
package main

import (
	"emsg-daemon/internal/config"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("EMSG_DATABASE_URL", "sqlite3://test.db")
	os.Setenv("EMSG_DOMAIN", "testdomain.com")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.DatabaseURL != "sqlite3://test.db" {
		t.Errorf("expected DatabaseURL to be 'sqlite3://test.db', got %s", cfg.DatabaseURL)
	}
	if cfg.Domain != "testdomain.com" {
		t.Errorf("expected Domain to be 'testdomain.com', got %s", cfg.Domain)
	}
}

func TestLoadConfigWWWDirAndLogFileDefaults(t *testing.T) {
	os.Unsetenv("EMSG_WWW_DIR")
	os.Unsetenv("EMSG_LOG_FILE")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(cfg.WWWDir, "www") {
		t.Errorf("expected WWWDir to end with 'www', got %s", cfg.WWWDir)
	}
	if !strings.HasSuffix(cfg.LogFile, "emsg.log") {
		t.Errorf("expected LogFile to end with 'emsg.log', got %s", cfg.LogFile)
	}
}

func TestLoadConfigWWWDirAndLogFileFromEnv(t *testing.T) {
	os.Setenv("EMSG_WWW_DIR", "/custom/www")
	os.Setenv("EMSG_LOG_FILE", "/custom/emsg.log")
	defer os.Unsetenv("EMSG_WWW_DIR")
	defer os.Unsetenv("EMSG_LOG_FILE")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.WWWDir != "/custom/www" {
		t.Errorf("expected WWWDir '/custom/www', got %s", cfg.WWWDir)
	}
	if cfg.LogFile != "/custom/emsg.log" {
		t.Errorf("expected LogFile '/custom/emsg.log', got %s", cfg.LogFile)
	}
}

func TestLoadConfigFromFileNewKeys(t *testing.T) {
	content := "EMSG_WWW_DIR=/srv/emsg/www\nEMSG_LOG_FILE=/var/log/emsg.log\nEMSG_PORT=9000\n"
	tmp, err := os.CreateTemp("", "emsg-*.conf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.WriteString(content)
	tmp.Close()

	cfg, err := config.LoadConfigFromFile(tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.WWWDir != "/srv/emsg/www" {
		t.Errorf("expected WWWDir '/srv/emsg/www', got %s", cfg.WWWDir)
	}
	if cfg.LogFile != "/var/log/emsg.log" {
		t.Errorf("expected LogFile '/var/log/emsg.log', got %s", cfg.LogFile)
	}
	if cfg.Port != "9000" {
		t.Errorf("expected Port '9000', got %s", cfg.Port)
	}
}

func TestLoadConfigFromFileAbsent(t *testing.T) {
	_, err := config.LoadConfigFromFile(filepath.Join(t.TempDir(), "nonexistent.conf"))
	if err == nil {
		t.Error("expected error for absent config file, got nil")
	}
}

func TestLoadConfigFromFileMaxConnectionsStrconv(t *testing.T) {
	content := "EMSG_MAX_CONNECTIONS=42\n"
	tmp, err := os.CreateTemp("", "emsg-*.conf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.WriteString(content)
	tmp.Close()

	cfg, err := config.LoadConfigFromFile(tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxConnections != 42 {
		t.Errorf("expected MaxConnections 42, got %d", cfg.MaxConnections)
	}
}
