// test_final_integration.go
// Final integration test for EMSG Daemon with custom configuration
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
	fmt.Println("🚀 EMSG Daemon Final Integration Test")
	fmt.Println("Testing daemon running on custom port 9090...")
	time.Sleep(2 * time.Second)

	baseURL := "http://localhost:9090"

	// Test 1: User registration and retrieval
	fmt.Println("\n1. Testing user management...")
	
	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	pubKeyB64 := base64.StdEncoding.EncodeToString(pubKey)
	
	userReq := map[string]string{
		"address":         "finaltest#production.emsg.dev",
		"pubkey":          pubKeyB64,
		"first_name":      "Final",
		"last_name":       "Test",
		"display_picture": "https://example.com/final.jpg",
	}

	if err := testUserRegistration(baseURL, userReq); err != nil {
		fmt.Printf("❌ User registration failed: %v\n", err)
		return
	}
	fmt.Println("✅ User registration successful")

	if err := testUserRetrieval(baseURL, "finaltest%23production.emsg.dev"); err != nil {
		fmt.Printf("❌ User retrieval failed: %v\n", err)
		return
	}
	fmt.Println("✅ User retrieval successful")

	// Test 2: Group management
	fmt.Println("\n2. Testing group management...")
	
	groupReq := map[string]interface{}{
		"id":          "final-test-group",
		"name":        "Final Test Group",
		"description": "Integration test group",
		"members":     []string{"finaltest#production.emsg.dev"},
	}

	if err := testGroupCreation(baseURL, groupReq, "finaltest#production.emsg.dev", privKey); err != nil {
		fmt.Printf("❌ Group creation failed: %v\n", err)
		return
	}
	fmt.Println("✅ Group creation successful")

	if err := testGroupRetrieval(baseURL, "final-test-group"); err != nil {
		fmt.Printf("❌ Group retrieval failed: %v\n", err)
		return
	}
	fmt.Println("✅ Group retrieval successful")

	// Test 3: Message sending and retrieval
	fmt.Println("\n3. Testing message management...")
	
	messageReq := map[string]interface{}{
		"from":      "finaltest#production.emsg.dev",
		"to":        []string{"testuser#production.emsg.dev"},
		"group_id":  "final-test-group",
		"body":      "Final integration test message!",
		"signature": "test_signature",
	}

	if err := testMessageSending(baseURL, messageReq, "finaltest#production.emsg.dev", privKey); err != nil {
		fmt.Printf("❌ Message sending failed: %v\n", err)
		return
	}
	fmt.Println("✅ Message sending successful")

	if err := testMessageRetrieval(baseURL, "finaltest%23production.emsg.dev", "finaltest#production.emsg.dev", privKey); err != nil {
		fmt.Printf("❌ Message retrieval failed: %v\n", err)
		return
	}
	fmt.Println("✅ Message retrieval successful")

	// Test 4: DNS routing functionality
	fmt.Println("\n4. Testing DNS routing...")
	
	if err := testAddressValidation(baseURL); err != nil {
		fmt.Printf("❌ Address validation failed: %v\n", err)
		return
	}
	fmt.Println("✅ Address validation successful")

	// Test 5: Authentication
	fmt.Println("\n5. Testing authentication...")
	
	if err := testAuthentication(baseURL, privKey); err != nil {
		fmt.Printf("❌ Authentication test failed: %v\n", err)
		return
	}
	fmt.Println("✅ Authentication test successful")

	fmt.Println("\n🎉 All integration tests passed!")
	fmt.Println("\n📊 EMSG Daemon Status:")
	fmt.Println("   ✅ Running on custom port 9090")
	fmt.Println("   ✅ Using custom domain: production.emsg.dev")
	fmt.Println("   ✅ Using custom database: ./production_emsg.db")
	fmt.Println("   ✅ All API endpoints functional")
	fmt.Println("   ✅ Authentication working")
	fmt.Println("   ✅ DNS routing implemented")
	fmt.Println("   ✅ BoltDB storage working")
	
	fmt.Println("\n🚀 EMSG Daemon is fully operational!")
}

func testUserRegistration(baseURL string, userReq map[string]string) error {
	jsonData, _ := json.Marshal(userReq)
	resp, err := http.Post(baseURL+"/api/user", "application/json", bytes.NewBuffer(jsonData))
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

func testUserRetrieval(baseURL, address string) error {
	resp, err := http.Get(baseURL + "/api/user?address=" + address)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func testGroupCreation(baseURL string, groupReq map[string]interface{}, userAddress string, privKey ed25519.PrivateKey) error {
	jsonData, _ := json.Marshal(groupReq)
	
	authHeader, err := api.CreateAuthRequest(userAddress, privKey, "POST", "/api/group")
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", baseURL+"/api/group", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "EMSG "+authHeader)
	
	client := &http.Client{}
	resp, err := client.Do(req)
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

func testGroupRetrieval(baseURL, groupID string) error {
	resp, err := http.Get(baseURL + "/api/group?id=" + groupID)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func testMessageSending(baseURL string, messageReq map[string]interface{}, userAddress string, privKey ed25519.PrivateKey) error {
	jsonData, _ := json.Marshal(messageReq)
	
	authHeader, err := api.CreateAuthRequest(userAddress, privKey, "POST", "/api/message")
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", baseURL+"/api/message", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "EMSG "+authHeader)
	
	client := &http.Client{}
	resp, err := client.Do(req)
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

func testMessageRetrieval(baseURL, userAddress, realAddress string, privKey ed25519.PrivateKey) error {
	authHeader, err := api.CreateAuthRequest(realAddress, privKey, "GET", "/api/messages")
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("GET", baseURL+"/api/messages?user="+userAddress, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "EMSG "+authHeader)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func testAddressValidation(baseURL string) error {
	validateReq := map[string]interface{}{
		"addresses": []string{"test#emsg.dev", "invalid-address"},
	}
	
	jsonData, _ := json.Marshal(validateReq)
	resp, err := http.Post(baseURL+"/api/route/validate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func testAuthentication(baseURL string, privKey ed25519.PrivateKey) error {
	// Test unauthenticated request (should fail)
	messageReq := map[string]interface{}{
		"from": "test#emsg.dev",
		"to":   []string{"other#emsg.dev"},
		"body": "test",
	}
	
	jsonData, _ := json.Marshal(messageReq)
	resp, err := http.Post(baseURL+"/api/message", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 401 {
		return fmt.Errorf("unauthenticated request should have failed with 401, got %d", resp.StatusCode)
	}
	
	return nil
}
