// auth_test.go
// Tests for public key registration and signature verification
package main

import (
	"crypto/ed25519"
	"emsg-daemon/internal/auth"
	"encoding/base64"
	"testing"
)

func TestRegisterUserAndVerifySignature(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	pubB64 := base64.StdEncoding.EncodeToString(pub)
	user, err := auth.RegisterUser("alice#emsg.dev", pubB64, "Alice", "B.", "Smith", "http://img/alice.png")
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}
	msg := []byte("hello world")
	sig := ed25519.Sign(priv, msg)
	if !auth.VerifySignature(user.PubKey, msg, sig) {
		t.Error("signature verification failed")
	}
}
