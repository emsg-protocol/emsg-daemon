// auth.go
// Key registration & signature verification for EMSG Daemon
package auth

import (
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
)

type User struct {
	Address        string            `json:"address"`
	PubKey         ed25519.PublicKey `json:"-"` // Handle separately
	FirstName      string            `json:"first_name"`
	MiddleName     string            `json:"middle_name"`
	LastName       string            `json:"last_name"`
	DisplayPicture string            `json:"display_picture"`
}

// MarshalJSON custom JSON marshaling for User
func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		PubKey string `json:"pubkey"`
		*Alias
	}{
		PubKey: base64.StdEncoding.EncodeToString(u.PubKey),
		Alias:  (*Alias)(&u),
	})
}

// UnmarshalJSON custom JSON unmarshaling for User
func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := &struct {
		PubKey string `json:"pubkey"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	pubKey, err := base64.StdEncoding.DecodeString(aux.PubKey)
	if err != nil {
		return err
	}
	u.PubKey = ed25519.PublicKey(pubKey)
	return nil
}

// RegisterUser registers a user's public key and profile fields
func RegisterUser(address, pubKeyBase64, firstName, middleName, lastName, displayPicture string) (*User, error) {
	pubKey, err := base64.StdEncoding.DecodeString(pubKeyBase64)
	if err != nil {
		return nil, err
	}
	if len(pubKey) != ed25519.PublicKeySize {
		return nil, errors.New("invalid public key size")
	}
	return &User{
		Address:        address,
		PubKey:         ed25519.PublicKey(pubKey),
		FirstName:      firstName,
		MiddleName:     middleName,
		LastName:       lastName,
		DisplayPicture: displayPicture,
	}, nil
}

// VerifySignature verifies a message signature
func VerifySignature(pubKey ed25519.PublicKey, message, sig []byte) bool {
	return ed25519.Verify(pubKey, message, sig)
}

// StoreUser inserts a user and their public key/profile fields into the database
func StoreUser(db *sql.DB, user *User) error {
	pubKeyB64 := base64.StdEncoding.EncodeToString(user.PubKey)
	_, err := db.Exec(`INSERT OR REPLACE INTO users (address, pubkey, first_name, middle_name, last_name, display_picture) VALUES (?, ?, ?, ?, ?, ?)`,
		user.Address, pubKeyB64, user.FirstName, user.MiddleName, user.LastName, user.DisplayPicture)
	return err
}

// GetUser retrieves a user and their public key/profile fields from the database
func GetUser(db *sql.DB, address string) (*User, error) {
	row := db.QueryRow(`SELECT address, pubkey, first_name, middle_name, last_name, display_picture FROM users WHERE address = ?`, address)
	var addr, pubKeyB64, firstName, middleName, lastName, displayPicture string
	if err := row.Scan(&addr, &pubKeyB64, &firstName, &middleName, &lastName, &displayPicture); err != nil {
		return nil, err
	}
	pubKey, err := base64.StdEncoding.DecodeString(pubKeyB64)
	if err != nil {
		return nil, err
	}
	return &User{
		Address:        addr,
		PubKey:         ed25519.PublicKey(pubKey),
		FirstName:      firstName,
		MiddleName:     middleName,
		LastName:       lastName,
		DisplayPicture: displayPicture,
	}, nil
}
