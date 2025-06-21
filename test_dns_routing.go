// test_dns_routing.go
// Test DNS routing functionality
package main

import (
	"fmt"
	"emsg-daemon/internal/router"
)

func main() {
	fmt.Println("Testing DNS Routing Functionality...")

	// Test 1: Address validation
	fmt.Println("\n1. Testing address validation...")
	
	validAddresses := []string{
		"alice#emsg.dev",
		"bob#example.com",
		"user123#test.domain.org",
	}
	
	invalidAddresses := []string{
		"alice",           // missing domain
		"#emsg.dev",      // missing user
		"alice@emsg.dev", // wrong separator
		"alice#",         // missing domain
		"alice#nodot",    // domain without dot
	}
	
	fmt.Println("Valid addresses:")
	for _, addr := range validAddresses {
		if err := router.ValidateAddress(addr); err != nil {
			fmt.Printf("❌ %s: %v\n", addr, err)
		} else {
			fmt.Printf("✅ %s: valid\n", addr)
		}
	}
	
	fmt.Println("\nInvalid addresses:")
	for _, addr := range invalidAddresses {
		if err := router.ValidateAddress(addr); err != nil {
			fmt.Printf("✅ %s: correctly rejected - %v\n", addr, err)
		} else {
			fmt.Printf("❌ %s: should have been rejected\n", addr)
		}
	}

	// Test 2: Route info parsing
	fmt.Println("\n2. Testing route info parsing...")
	
	// Test JSON format
	jsonRoute := `{"server":"https://emsg.example.com:8080","pubkey":"abc123","version":"1.0","ttl":3600}`
	routeInfo, err := router.ParseRouteInfo(jsonRoute)
	if err != nil {
		fmt.Printf("❌ Failed to parse JSON route: %v\n", err)
	} else {
		fmt.Printf("✅ JSON route parsed: server=%s, version=%s, ttl=%d\n", 
			routeInfo.Server, routeInfo.Version, routeInfo.TTL)
	}
	
	// Test simple URL format
	urlRoute := "https://emsg.example.com:8080"
	routeInfo, err = router.ParseRouteInfo(urlRoute)
	if err != nil {
		fmt.Printf("❌ Failed to parse URL route: %v\n", err)
	} else {
		fmt.Printf("✅ URL route parsed: server=%s, version=%s, ttl=%d\n", 
			routeInfo.Server, routeInfo.Version, routeInfo.TTL)
	}

	// Test 3: Message routing (simulation)
	fmt.Println("\n3. Testing message routing logic...")
	
	// This will fail with DNS lookup errors, but we can test the logic
	recipients := []string{"alice#emsg.dev", "bob#example.com"}
	
	fmt.Printf("Attempting to route message to: %v\n", recipients)
	routes, err := router.RouteMessage(recipients)
	if err != nil {
		fmt.Printf("⚠️  Expected DNS lookup failure: %v\n", err)
		fmt.Println("   (This is normal - we don't have real DNS records set up)")
	} else {
		fmt.Printf("✅ Routes determined: %v\n", routes)
	}

	// Test 4: Local domain checking
	fmt.Println("\n4. Testing local domain checking...")
	
	localDomains := []string{"emsg.dev", "localhost", "example.local"}
	testDomains := []string{"emsg.dev", "example.com", "localhost", "google.com"}
	
	for _, domain := range testDomains {
		isLocal := router.IsLocalDomain(domain, localDomains)
		if isLocal {
			fmt.Printf("✅ %s: local domain\n", domain)
		} else {
			fmt.Printf("🌐 %s: remote domain\n", domain)
		}
	}

	// Test 5: DNS lookup (will likely fail but shows the functionality)
	fmt.Println("\n5. Testing DNS lookup...")
	
	testAddress := "test#emsg.dev"
	fmt.Printf("Looking up route for: %s\n", testAddress)
	route, err := router.LookupRoute(testAddress)
	if err != nil {
		fmt.Printf("⚠️  DNS lookup failed (expected): %v\n", err)
		fmt.Println("   To test this properly, set up a TXT record at _emsg.emsg.dev")
	} else {
		fmt.Printf("✅ Route found: %s\n", route)
	}

	fmt.Println("\n🎉 DNS Routing functionality testing completed!")
	fmt.Println("\nNote: DNS lookups will fail without real TXT records.")
	fmt.Println("To test fully, create TXT records like:")
	fmt.Println("  _emsg.yourdomain.com TXT \"https://your-emsg-server.com:8080\"")
	fmt.Println("  or")
	fmt.Println("  _emsg.yourdomain.com TXT \"{\\\"server\\\":\\\"https://your-emsg-server.com:8080\\\",\\\"version\\\":\\\"1.0\\\"}\"")
}
