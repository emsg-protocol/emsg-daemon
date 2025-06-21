// send_emsg_message.go
// Send EMSG messages with proper authentication
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

	"emsg-daemon/api"
)

func main() {
	fmt.Println("📧 EMSG Message Sender")
	fmt.Println("======================")
	
	// Configuration
	emsgServer := "https://emsg.sandipwalke.com"
	fromAddress := "sandip#sandipwalke.com"
	
	// Step 1: Generate or load your Ed25519 keys
	fmt.Println("\n1. Setting up cryptographic keys...")
	
	// For demo, we'll generate new keys
	// In practice, you'd load your saved private key
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		fmt.Printf("❌ Failed to generate keys: %v\n", err)
		return
	}
	
	pubKeyB64 := base64.StdEncoding.EncodeToString(pubKey)
	fmt.Printf("✅ Keys ready (Public: %s...)\n", pubKeyB64[:16])
	
	// Step 2: Register sender (if not already registered)
	fmt.Println("\n2. Ensuring sender is registered...")
	if err := registerUser(emsgServer, fromAddress, pubKeyB64); err != nil {
		fmt.Printf("⚠️  Registration: %v\n", err)
		// Continue anyway - user might already be registered
	} else {
		fmt.Println("✅ Sender registered successfully")
	}
	
	// Step 3: Create and send a message
	fmt.Println("\n3. Creating message...")
	
	message := map[string]interface{}{
		"from":      fromAddress,
		"to":        []string{"recipient#example.com"}, // Change this to real recipient
		"cc":        []string{},
		"group_id":  "",
		"body":      "Hello! This is a test message from the EMSG protocol. 🚀",
		"signature": "will_be_generated",
	}
	
	fmt.Printf("   From: %s\n", message["from"])
	fmt.Printf("   To: %v\n", message["to"])
	fmt.Printf("   Body: %s\n", message["body"])
	
	// Step 4: Send the message with authentication
	fmt.Println("\n4. Sending message with authentication...")
	
	if err := sendAuthenticatedMessage(emsgServer, message, fromAddress, privKey); err != nil {
		fmt.Printf("❌ Failed to send message: %v\n", err)
		return
	}
	
	fmt.Println("✅ Message sent successfully!")
	
	// Step 5: Demonstrate retrieving messages
	fmt.Println("\n5. Retrieving messages for sender...")
	
	messages, err := getMessages(emsgServer, fromAddress, privKey)
	if err != nil {
		fmt.Printf("❌ Failed to retrieve messages: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Retrieved %d messages\n", len(messages))
	for i, msg := range messages {
		if msgMap, ok := msg.(map[string]interface{}); ok {
			from := msgMap["from"]
			body := msgMap["body"]
			fmt.Printf("   Message %d: From %v - %v\n", i+1, from, body)
		}
	}
	
	fmt.Println("\n======================")
	fmt.Println("📧 Message Sending Complete!")
	fmt.Println("======================")
	
	fmt.Println("\n📋 How to send messages:")
	fmt.Println("1. Have your Ed25519 private key")
	fmt.Println("2. Create message with recipient addresses")
	fmt.Println("3. Sign the request with your private key")
	fmt.Println("4. Send POST request to /api/message")
	
	fmt.Println("\n🔐 Security Notes:")
	fmt.Println("- Messages require Ed25519 signature authentication")
	fmt.Println("- Your private key proves you own the sender address")
	fmt.Println("- Recipients are routed via DNS TXT record lookup")
	fmt.Println("- All communication uses HTTPS")
}

func registerUser(server, address, pubKey string) error {
	userReq := map[string]string{
		"address":         address,
		"pubkey":          pubKey,
		"first_name":      "Demo",
		"last_name":       "User",
		"display_picture": "",
	}
	
	jsonData, _ := json.Marshal(userReq)
	client := &http.Client{Timeout: 30 * time.Second}
	
	resp, err := client.Post(server+"/api/user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registration failed: %d - %s", resp.StatusCode, string(body))
	}
	
	return nil
}

func sendAuthenticatedMessage(server string, message map[string]interface{}, fromAddress string, privKey ed25519.PrivateKey) error {
	// Create authentication header
	authHeader, err := api.CreateAuthRequest(fromAddress, privKey, "POST", "/api/message")
	if err != nil {
		return fmt.Errorf("failed to create auth header: %w", err)
	}
	
	// Prepare message JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", server+"/api/message", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "EMSG "+authHeader)
	
	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		return fmt.Errorf("message sending failed: %d - %s", resp.StatusCode, string(body))
	}
	
	return nil
}

func getMessages(server, userAddress string, privKey ed25519.PrivateKey) ([]interface{}, error) {
	// Create authentication header
	authHeader, err := api.CreateAuthRequest(userAddress, privKey, "GET", "/api/messages")
	if err != nil {
		return nil, fmt.Errorf("failed to create auth header: %w", err)
	}
	
	// URL encode the address
	encodedAddress := userAddress
	if userAddress == "sandip#sandipwalke.com" {
		encodedAddress = "sandip%23sandipwalke.com"
	}
	
	// Create HTTP request
	req, err := http.NewRequest("GET", server+"/api/messages?user="+encodedAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "EMSG "+authHeader)
	
	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("message retrieval failed: %d - %s", resp.StatusCode, string(body))
	}
	
	// Parse messages
	var messages []interface{}
	if err := json.Unmarshal(body, &messages); err != nil {
		return nil, fmt.Errorf("failed to parse messages: %w", err)
	}
	
	return messages, nil
}
