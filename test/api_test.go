// api_test.go
// Tests for user registration and retrieval via REST API
package main

import (
	"bytes"
	"crypto/ed25519"
	"emsg-daemon/api"
	"emsg-daemon/internal/auth"
	"emsg-daemon/internal/storage"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestAPIRegisterAndGetUser(t *testing.T) {
	// Use a temp BoltDB file
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := storage.InitBoltDB(dbPath)
	if err != nil {
		t.Fatalf("InitBoltDB failed: %v", err)
	}
	defer db.Close()
	defer os.Remove(dbPath)

	apiHandler := &api.BoltAPI{DB: db}

	pub, _, _ := ed25519.GenerateKey(nil)
	pubB64 := base64.StdEncoding.EncodeToString(pub)
	userReq := map[string]string{
		"address":         "alice#emsg.dev",
		"pubkey":          pubB64,
		"first_name":      "Alice",
		"middle_name":     "B.",
		"last_name":       "Smith",
		"display_picture": "http://img/alice.png",
	}
	body, _ := json.Marshal(userReq)
	req := httptest.NewRequest("POST", "/api/user", bytes.NewReader(body))
	w := httptest.NewRecorder()
	apiHandler.ApiRegisterUser(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d. Response: %s", w.Code, w.Body.String())
	}

	getReq := httptest.NewRequest("GET", "/api/user?address=alice%23emsg.dev", nil)
	getW := httptest.NewRecorder()
	apiHandler.ApiGetUser(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", getW.Code)
	}
	var respUser auth.User
	if err := json.NewDecoder(getW.Body).Decode(&respUser); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if respUser.FirstName != "Alice" || respUser.LastName != "Smith" || respUser.DisplayPicture != "http://img/alice.png" {
		t.Error("user profile fields not returned correctly")
	}
}
