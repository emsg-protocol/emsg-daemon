// auth.go
// Key registration & signature verification for EMSG Daemon
package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
)

type User struct {
	Address   string
	PubKey    ed25519.PublicKey
}

// RegisterUser registers a user's public key
func RegisterUser(address string, pubKeyBase64 string) (*User, error) {
	pubKey, err := base64.StdEncoding.DecodeString(pubKeyBase64)
	if err != nil {
		return nil, err
	}
	if len(pubKey) != ed25519.PublicKeySize {
		return nil, errors.New("invalid public key size")
	}
	return &User{Address: address, PubKey: ed25519.PublicKey(pubKey)}, nil
}

// VerifySignature verifies a message signature
func VerifySignature(pubKey ed25519.PublicKey, message, sig []byte) bool {
	return ed25519.Verify(pubKey, message, sig)
}
