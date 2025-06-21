// simple_test.go
// Simple test for message and group endpoints
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Testing Message and Group endpoints...")
	time.Sleep(1 * time.Second)

	// Test 1: Create a group
	fmt.Println("\n1. Creating a group...")
	groupData := map[string]interface{}{
		"id":          "test-group",
		"name":        "Test Group",
		"description": "A test group for EMSG",
		"display_pic": "https://example.com/group.jpg",
		"members":     []string{"alice#emsg.dev", "bob#emsg.dev"},
	}

	jsonData, _ := json.Marshal(groupData)
	resp, err := http.Post("http://localhost:8080/api/group", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Error creating group: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 201 {
		fmt.Println("✅ Group created successfully!")
		fmt.Printf("Response: %s\n", string(body))
	} else {
		fmt.Printf("❌ Group creation failed: %d - %s\n", resp.StatusCode, string(body))
	}

	// Test 2: Get the group
	fmt.Println("\n2. Retrieving the group...")
	resp, err = http.Get("http://localhost:8080/api/group?id=test-group")
	if err != nil {
		fmt.Printf("❌ Error getting group: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Println("✅ Group retrieved successfully!")
		fmt.Printf("Response: %s\n", string(body))
	} else {
		fmt.Printf("❌ Group retrieval failed: %d - %s\n", resp.StatusCode, string(body))
	}

	// Test 3: Send a message
	fmt.Println("\n3. Sending a message...")
	messageData := map[string]interface{}{
		"from":      "alice#emsg.dev",
		"to":        []string{"bob#emsg.dev"},
		"cc":        []string{},
		"group_id":  "test-group",
		"body":      "Hello from Alice! This is a test message.",
		"signature": "dummy_signature_for_testing",
	}

	jsonData, _ = json.Marshal(messageData)
	resp, err = http.Post("http://localhost:8080/api/message", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Error sending message: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode == 201 {
		fmt.Println("✅ Message sent successfully!")
		fmt.Printf("Response: %s\n", string(body))
	} else {
		fmt.Printf("❌ Message sending failed: %d - %s\n", resp.StatusCode, string(body))
	}

	// Test 4: Get messages for Bob
	fmt.Println("\n4. Getting messages for Bob...")
	resp, err = http.Get("http://localhost:8080/api/messages?user=bob%23emsg.dev")
	if err != nil {
		fmt.Printf("❌ Error getting messages: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Println("✅ Messages retrieved successfully!")
		fmt.Printf("Response: %s\n", string(body))
	} else {
		fmt.Printf("❌ Message retrieval failed: %d - %s\n", resp.StatusCode, string(body))
	}

	fmt.Println("\n🎉 Testing completed!")
}
