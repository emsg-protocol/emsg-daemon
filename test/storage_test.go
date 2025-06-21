// storage_test.go
// Tests for database initialization and basic storage
package main

import (
	"emsg-daemon/internal/storage"
	"os"
	"testing"
)

func TestInitDB(t *testing.T) {
	dbFile := "test_emsg.db"
	defer os.Remove(dbFile)
	db, err := storage.InitDB(dbFile)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	if db == nil {
		t.Error("expected db instance, got nil")
	}
	db.Close()
}
