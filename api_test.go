// api_test.go
// Tests for user registration and retrieval via REST API
package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIRegisterAndGetUser(t *testing.T) {
	db, _ := InitDB(":memory:")
	InitSchema(db)
	api := &API{db: db}

	pub, _, _ := ed25519.GenerateKey(nil)
	pubB64 := base64.StdEncoding.EncodeToString(pub)
	userReq := map[string]string{
		"address": "alice#emsg.dev",
		"pubkey": pubB64,
		"first_name": "Alice",
		"middle_name": "B.",
		"last_name": "Smith",
		"display_picture": "http://img/alice.png",
	}
	body, _ := json.Marshal(userReq)
	req := httptest.NewRequest("POST", "/api/user", bytes.NewReader(body))
	w := httptest.NewRecorder()
	api.apiRegisterUser(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", w.Code)
	}

	getReq := httptest.NewRequest("GET", "/api/user?address=alice#emsg.dev", nil)
	getW := httptest.NewRecorder()
	api.apiGetUser(getW, getReq)
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
