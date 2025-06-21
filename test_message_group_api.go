// test_message_group_api.go
// Test script for EMSG Daemon Message and Group API endpoints
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

type Message struct {
	From      string   `json:"from"`
	To        []string `json:"to"`
	CC        []string `json:"cc"`
	GroupID   string   `json:"group_id"`
	Body      string   `json:"body"`
	Signature string   `json:"signature"`
}

type Group struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	DisplayPic  string   `json:"display_pic"`
	Members     []string `json:"members"`
}

func main() {
	fmt.Println("Testing EMSG Daemon Message and Group API...")
	time.Sleep(2 * time.Second)

	// Test 1: Register test users
	fmt.Println("\n1. Registering test users...")
	
	// Register Alice
	pubKey1, _, _ := ed25519.GenerateKey(nil)
	pubKeyB64_1 := base64.StdEncoding.EncodeToString(pubKey1)
	alice := UserRegistrationRequest{
		Address:        "alice#emsg.dev",
		PubKey:         pubKeyB64_1,
		FirstName:      "Alice",
		LastName:       "Smith",
		DisplayPicture: "https://example.com/alice.jpg",
	}
	
	if err := registerUser(alice); err != nil {
		fmt.Printf("❌ Failed to register Alice: %v\n", err)
		return
	}
	fmt.Println("✅ Alice registered successfully")

	// Register Bob
	pubKey2, _, _ := ed25519.GenerateKey(nil)
	pubKeyB64_2 := base64.StdEncoding.EncodeToString(pubKey2)
	bob := UserRegistrationRequest{
		Address:        "bob#emsg.dev",
		PubKey:         pubKeyB64_2,
		FirstName:      "Bob",
		LastName:       "Johnson",
		DisplayPicture: "https://example.com/bob.jpg",
	}
	
	if err := registerUser(bob); err != nil {
		fmt.Printf("❌ Failed to register Bob: %v\n", err)
		return
	}
	fmt.Println("✅ Bob registered successfully")

	// Test 2: Create a group
	fmt.Println("\n2. Testing group creation...")
	group := Group{
		ID:          "dev-team",
		Name:        "Development Team",
		Description: "EMSG Development Team Chat",
		DisplayPic:  "https://example.com/dev-team.jpg",
		Members:     []string{"alice#emsg.dev", "bob#emsg.dev"},
	}

	if err := createGroup(group); err != nil {
		fmt.Printf("❌ Failed to create group: %v\n", err)
		return
	}
	fmt.Println("✅ Group created successfully")

	// Test 3: Retrieve group info
	fmt.Println("\n3. Testing group retrieval...")
	retrievedGroup, err := getGroup("dev-team")
	if err != nil {
		fmt.Printf("❌ Failed to retrieve group: %v\n", err)
		return
	}
	fmt.Printf("✅ Group retrieved: %s (%d members)\n", retrievedGroup.Name, len(retrievedGroup.Members))

	// Test 4: Send a message
	fmt.Println("\n4. Testing message sending...")
	message := Message{
		From:      "alice#emsg.dev",
		To:        []string{"bob#emsg.dev"},
		CC:        []string{},
		GroupID:   "dev-team",
		Body:      "Hello Bob! Welcome to the EMSG development team!",
		Signature: "dummy_signature",
	}

	if err := sendMessage(message); err != nil {
		fmt.Printf("❌ Failed to send message: %v\n", err)
		return
	}
	fmt.Println("✅ Message sent successfully")

	// Test 5: Send a group message
	fmt.Println("\n5. Testing group message...")
	groupMessage := Message{
		From:      "bob#emsg.dev",
		To:        []string{"alice#emsg.dev"},
		CC:        []string{},
		GroupID:   "dev-team",
		Body:      "Thanks Alice! Excited to work on EMSG together!",
		Signature: "dummy_signature",
	}

	if err := sendMessage(groupMessage); err != nil {
		fmt.Printf("❌ Failed to send group message: %v\n", err)
		return
	}
	fmt.Println("✅ Group message sent successfully")

	// Test 6: Retrieve messages for Alice
	fmt.Println("\n6. Testing message retrieval for Alice...")
	aliceMessages, err := getMessages("alice#emsg.dev")
	if err != nil {
		fmt.Printf("❌ Failed to retrieve Alice's messages: %v\n", err)
		return
	}
	fmt.Printf("✅ Retrieved %d messages for Alice\n", len(aliceMessages))
	for i, msg := range aliceMessages {
		fmt.Printf("   Message %d: From %s - %s\n", i+1, msg.From, msg.Body[:min(50, len(msg.Body))]+"...")
	}

	// Test 7: Retrieve messages for Bob
	fmt.Println("\n7. Testing message retrieval for Bob...")
	bobMessages, err := getMessages("bob#emsg.dev")
	if err != nil {
		fmt.Printf("❌ Failed to retrieve Bob's messages: %v\n", err)
		return
	}
	fmt.Printf("✅ Retrieved %d messages for Bob\n", len(bobMessages))
	for i, msg := range bobMessages {
		fmt.Printf("   Message %d: From %s - %s\n", i+1, msg.From, msg.Body[:min(50, len(msg.Body))]+"...")
	}

	fmt.Println("\n🎉 Message and Group API testing completed!")
}

func registerUser(user UserRegistrationRequest) error {
	jsonData, _ := json.Marshal(user)
	resp, err := http.Post("http://localhost:8080/api/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func createGroup(group Group) error {
	jsonData, _ := json.Marshal(group)
	resp, err := http.Post("http://localhost:8080/api/group", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func getGroup(id string) (*Group, error) {
	resp, err := http.Get("http://localhost:8080/api/group?id=" + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	
	var group Group
	if err := json.NewDecoder(resp.Body).Decode(&group); err != nil {
		return nil, err
	}
	return &group, nil
}

func sendMessage(message Message) error {
	jsonData, _ := json.Marshal(message)
	resp, err := http.Post("http://localhost:8080/api/message", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func getMessages(user string) ([]Message, error) {
	// URL encode the user address
	encodedUser := user
	if user == "alice#emsg.dev" {
		encodedUser = "alice%23emsg.dev"
	} else if user == "bob#emsg.dev" {
		encodedUser = "bob%23emsg.dev"
	}
	
	resp, err := http.Get("http://localhost:8080/api/messages?user=" + encodedUser)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	
	var messages []Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
