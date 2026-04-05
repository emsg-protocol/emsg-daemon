//go:build ignore
// test_final_routing.go
// Final test of EMSG routing for sandipwalke.com
package main

import (
	"fmt"
	"emsg-daemon/internal/router"
)

func main() {
	fmt.Println("🔍 Final EMSG Routing Test for sandipwalke.com")
	fmt.Println("==============================================")
	
	testAddress := "sandip#sandipwalke.com"
	
	// Test 1: Address validation
	fmt.Println("\n1. Testing address validation...")
	if err := router.ValidateAddress(testAddress); err != nil {
		fmt.Printf("❌ Address validation failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Address valid: %s\n", testAddress)
	
	// Test 2: DNS route lookup
	fmt.Println("\n2. Testing DNS route lookup...")
	route, err := router.LookupRoute(testAddress)
	if err != nil {
		fmt.Printf("❌ Route lookup failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Route found: %s\n", route)
	
	// Test 3: Route parsing
	fmt.Println("\n3. Testing route parsing...")
	routeInfo, err := router.ParseRouteInfo(route)
	if err != nil {
		fmt.Printf("   → Simple URL format: %s\n", route)
	} else {
		fmt.Printf("✅ Structured route info:\n")
		fmt.Printf("   Server: %s\n", routeInfo.Server)
		fmt.Printf("   Version: %s\n", routeInfo.Version)
		fmt.Printf("   TTL: %d seconds\n", routeInfo.TTL)
	}
	
	// Test 4: Message routing
	fmt.Println("\n4. Testing message routing...")
	routes, err := router.RouteMessage([]string{testAddress})
	if err != nil {
		fmt.Printf("❌ Message routing failed: %v\n", err)
		return
	}
	
	fmt.Println("✅ Message routing successful!")
	for server, recipients := range routes {
		fmt.Printf("   Server: %s\n", server)
		fmt.Printf("   Recipients: %v\n", recipients)
	}
	
	// Test 5: Multiple addresses
	fmt.Println("\n5. Testing multiple address routing...")
	multiAddresses := []string{
		"sandip#sandipwalke.com",
		"admin#sandipwalke.com",
		"contact#sandipwalke.com",
	}
	
	multiRoutes, err := router.RouteMessage(multiAddresses)
	if err != nil {
		fmt.Printf("❌ Multi-address routing failed: %v\n", err)
		return
	}
	
	fmt.Println("✅ Multi-address routing successful!")
	for server, recipients := range multiRoutes {
		fmt.Printf("   Server: %s\n", server)
		fmt.Printf("   Recipients: %v\n", recipients)
	}
	
	fmt.Println("\n==============================================")
	fmt.Println("🎉 ALL ROUTING TESTS PASSED!")
	fmt.Println("==============================================")
	
	fmt.Println("\n📊 Final Status:")
	fmt.Println("✅ DNS TXT Record: Working")
	fmt.Println("✅ Route Discovery: Working")
	fmt.Println("✅ Message Routing: Working")
	fmt.Println("✅ Multi-user Support: Working")
	
	fmt.Println("\n🌐 Your domain is fully integrated into the EMSG network!")
	fmt.Println("Other EMSG users can now send messages to:")
	for _, addr := range multiAddresses {
		fmt.Printf("   - %s\n", addr)
	}
	
	fmt.Println("\n🔧 Optional DNS Optimization:")
	fmt.Println("Update your TXT record from:")
	fmt.Println("   https://emsg.sandipwalke.com:8765")
	fmt.Println("To:")
	fmt.Println("   https://emsg.sandipwalke.com")
	fmt.Println("(Remove :8765 since HTTPS uses port 443)")
}
