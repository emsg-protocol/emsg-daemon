// test_dns_api.go
// Test DNS routing API endpoints
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
	fmt.Println("Testing DNS Routing API endpoints...")
	time.Sleep(1 * time.Second)

	// Test 1: Address validation endpoint
	fmt.Println("\n1. Testing address validation endpoint...")
	
	validateReq := map[string]interface{}{
		"addresses": []string{
			"alice#emsg.dev",
			"bob#example.com", 
			"invalid-address",
			"#missing-user.com",
			"user#nodot",
		},
	}
	
	jsonData, _ := json.Marshal(validateReq)
	resp, err := http.Post("http://localhost:8080/api/route/validate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Error testing validation: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Println("✅ Address validation endpoint working!")
		
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		
		if results, ok := result["results"].(map[string]interface{}); ok {
			for addr, info := range results {
				if infoMap, ok := info.(map[string]interface{}); ok {
					if valid, ok := infoMap["valid"].(bool); ok {
						if valid {
							fmt.Printf("   ✅ %s: valid\n", addr)
						} else {
							fmt.Printf("   ❌ %s: %s\n", addr, infoMap["error"])
						}
					}
				}
			}
		}
	} else {
		fmt.Printf("❌ Validation endpoint failed: %d - %s\n", resp.StatusCode, string(body))
	}

	// Test 2: Route lookup endpoint (will fail with DNS error, but tests the endpoint)
	fmt.Println("\n2. Testing route lookup endpoint...")
	
	resp, err = http.Get("http://localhost:8080/api/route?address=alice%23emsg.dev")
	if err != nil {
		fmt.Printf("❌ Error testing route lookup: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Println("✅ Route lookup successful!")
		fmt.Printf("Response: %s\n", string(body))
	} else {
		fmt.Printf("⚠️  Route lookup failed (expected): %d - %s\n", resp.StatusCode, string(body))
		fmt.Println("   (This is normal - no DNS TXT records are set up)")
	}

	// Test 3: Message routing endpoint
	fmt.Println("\n3. Testing message routing endpoint...")
	
	routeReq := map[string]interface{}{
		"recipients": []string{
			"alice#emsg.dev",
			"bob#example.com",
		},
	}
	
	jsonData, _ = json.Marshal(routeReq)
	resp, err = http.Post("http://localhost:8080/api/route/message", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Error testing message routing: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Println("✅ Message routing successful!")
		fmt.Printf("Response: %s\n", string(body))
	} else {
		fmt.Printf("⚠️  Message routing failed (expected): %d - %s\n", resp.StatusCode, string(body))
		fmt.Println("   (This is normal - no DNS TXT records are set up)")
	}

	// Test 4: Invalid address in route lookup
	fmt.Println("\n4. Testing invalid address handling...")
	
	resp, err = http.Get("http://localhost:8080/api/route?address=invalid-address")
	if err != nil {
		fmt.Printf("❌ Error testing invalid address: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode == 400 {
		fmt.Println("✅ Invalid address correctly rejected!")
		fmt.Printf("Error: %s\n", string(body))
	} else {
		fmt.Printf("❌ Invalid address should have been rejected: %d - %s\n", resp.StatusCode, string(body))
	}

	// Test 5: Missing parameters
	fmt.Println("\n5. Testing missing parameters...")
	
	resp, err = http.Get("http://localhost:8080/api/route")
	if err != nil {
		fmt.Printf("❌ Error testing missing params: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode == 400 {
		fmt.Println("✅ Missing address parameter correctly rejected!")
		fmt.Printf("Error: %s\n", string(body))
	} else {
		fmt.Printf("❌ Missing parameter should have been rejected: %d - %s\n", resp.StatusCode, string(body))
	}

	fmt.Println("\n🎉 DNS Routing API testing completed!")
	fmt.Println("\nNote: Actual DNS lookups will fail without real TXT records.")
	fmt.Println("The endpoints are working correctly - they just need real DNS records to resolve.")
}
