// debug_storage.go
// Debug script to test BoltDB storage directly
package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"log"

	"emsg-daemon/internal/auth"
	"emsg-daemon/internal/storage"
)

func main() {
	fmt.Println("Testing BoltDB storage directly...")

	// Initialize BoltDB (same file as daemon)
	db, err := storage.InitBoltDB("./emsg.db")
	if err != nil {
		log.Fatalf("Failed to init BoltDB: %v", err)
	}
	defer db.Close()

	// Generate a test user
	pubKey, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}
	pubKeyB64 := base64.StdEncoding.EncodeToString(pubKey)

	user, err := auth.RegisterUser("alice#emsg.dev", pubKeyB64, "Alice", "B.", "Smith", "https://example.com/alice.jpg")
	if err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}

	fmt.Printf("Created user: %+v\n", user)

	// Store the user
	fmt.Println("Storing user...")
	err = storage.StoreUserBolt(db, user)
	if err != nil {
		log.Fatalf("Failed to store user: %v", err)
	}
	fmt.Println("User stored successfully!")

	// Retrieve the user
	fmt.Println("Retrieving user...")
	retrievedUser, err := storage.GetUserBolt(db, "alice#emsg.dev")
	if err != nil {
		log.Fatalf("Failed to retrieve user: %v", err)
	}

	fmt.Printf("Retrieved user: %+v\n", retrievedUser)
	fmt.Println("✅ Storage test successful!")
}
