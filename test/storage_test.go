// storage_test.go
// Tests for database initialization and basic storage
package main

import (
	"emsg-daemon/internal/storage"
	"path/filepath"
	"testing"
)

func TestInitDB(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test_emsg.db")
	db, err := storage.InitBoltDB(dbPath)
	if err != nil {
		t.Fatalf("InitBoltDB failed: %v", err)
	}
	if db == nil {
		t.Error("expected db instance, got nil")
	}
	db.Close()
}
