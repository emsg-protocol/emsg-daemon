// api_test.go
// Tests for user registration and retrieval via REST API
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIRegisterAndGetUser(t *testing.T) {
	// Setup in-memory DB and handler context as needed
	// This is a simplified example; in production, use a test DB and proper setup/teardown
	// Assume db is initialized and available globally

	// Register user
	userReq := map[string]string{
		"address": "alice#emsg.dev",
		"pubkey": "dGVzdHB1YmtleQ==", // base64 for 'testpubkey' (not a real Ed25519 key)
		"first_name": "Alice",
		"middle_name": "B.",
		"last_name": "Smith",
		"display_picture": "http://img/alice.png",
	}
	body, _ := json.Marshal(userReq)
	req := httptest.NewRequest("POST", "/api/user", bytes.NewReader(body))
	w := httptest.NewRecorder()
	apiRegisterUser(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", w.Code)
	}

	// Retrieve user
	getReq := httptest.NewRequest("GET", "/api/user?address=alice#emsg.dev", nil)
	getW := httptest.NewRecorder()
	apiGetUser(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", getW.Code)
	}
	var respUser User
	if err := json.NewDecoder(getW.Body).Decode(&respUser); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if respUser.FirstName != "Alice" || respUser.LastName != "Smith" || respUser.DisplayPicture != "http://img/alice.png" {
		t.Error("user profile fields not returned correctly")
	}
}
