// test_api.go
// Test script for EMSG Daemon API endpoints
package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type UserRegistrationRequest struct {
	Address        string `json:"address"`
	PubKey         string `json:"pubkey"`
	FirstName      string `json:"first_name"`
	MiddleName     string `json:"middle_name"`
	LastName       string `json:"last_name"`
	DisplayPicture string `json:"display_picture"`
}

type User struct {
	Address        string `json:"address"`
	FirstName      string `json:"first_name"`
	MiddleName     string `json:"middle_name"`
	LastName       string `json:"last_name"`
	DisplayPicture string `json:"display_picture"`
}

func main() {
	fmt.Println("Testing EMSG Daemon API...")

	// Wait a moment for the server to be ready
	time.Sleep(2 * time.Second)

	// Test 1: Generate Ed25519 key pair
	fmt.Println("\n1. Generating Ed25519 key pair...")
	pubKey, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		fmt.Printf("Failed to generate key: %v\n", err)
		return
	}
	pubKeyB64 := base64.StdEncoding.EncodeToString(pubKey)
	fmt.Printf("Generated public key: %s\n", pubKeyB64[:32]+"...")

	// Test 2: Register a new user
	fmt.Println("\n2. Testing user registration...")
	userReq := UserRegistrationRequest{
		Address:        "alice#emsg.dev",
		PubKey:         pubKeyB64,
		FirstName:      "Alice",
		MiddleName:     "B.",
		LastName:       "Smith",
		DisplayPicture: "https://example.com/alice.jpg",
	}

	jsonData, err := json.Marshal(userReq)
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/api/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to register user: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		fmt.Println("✅ User registration successful!")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ User registration failed: %d - %s\n", resp.StatusCode, string(body))
		return
	}

	// Test 3: Retrieve the registered user
	fmt.Println("\n3. Testing user retrieval...")
	// URL encode the address to handle the # character
	encodedAddress := "alice%23emsg.dev"
	resp, err = http.Get("http://localhost:8080/api/user?address=" + encodedAddress)
	if err != nil {
		fmt.Printf("Failed to get user: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var user User
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &user); err != nil {
			fmt.Printf("Failed to parse user response: %v\n", err)
			return
		}
		fmt.Println("✅ User retrieval successful!")
		fmt.Printf("   Address: %s\n", user.Address)
		fmt.Printf("   Name: %s %s %s\n", user.FirstName, user.MiddleName, user.LastName)
		fmt.Printf("   Picture: %s\n", user.DisplayPicture)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ User retrieval failed: %d - %s\n", resp.StatusCode, string(body))
		return
	}

	// Test 4: Try to register the same user again (should handle duplicates)
	fmt.Println("\n4. Testing duplicate user registration...")
	resp, err = http.Post("http://localhost:8080/api/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to test duplicate registration: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		fmt.Println("✅ Duplicate registration handled (overwrote existing user)")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("ℹ️  Duplicate registration response: %d - %s\n", resp.StatusCode, string(body))
	}

	// Test 5: Register another user
	fmt.Println("\n5. Testing second user registration...")
	pubKey2, _, _ := ed25519.GenerateKey(nil)
	pubKeyB64_2 := base64.StdEncoding.EncodeToString(pubKey2)

	userReq2 := UserRegistrationRequest{
		Address:        "bob#emsg.dev",
		PubKey:         pubKeyB64_2,
		FirstName:      "Bob",
		MiddleName:     "",
		LastName:       "Johnson",
		DisplayPicture: "https://example.com/bob.jpg",
	}

	jsonData2, _ := json.Marshal(userReq2)
	resp, err = http.Post("http://localhost:8080/api/user", "application/json", bytes.NewBuffer(jsonData2))
	if err != nil {
		fmt.Printf("Failed to register second user: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		fmt.Println("✅ Second user registration successful!")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Second user registration failed: %d - %s\n", resp.StatusCode, string(body))
	}

	fmt.Println("\n🎉 API testing completed!")
}
