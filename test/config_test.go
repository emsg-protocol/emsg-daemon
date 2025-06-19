// config_test.go
// Tests for config loading in EMSG Daemon
package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("EMSG_DATABASE_URL", "sqlite3://test.db")
	os.Setenv("EMSG_DOMAIN", "testdomain.com")
	cfg, err := LoadConfig()
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
