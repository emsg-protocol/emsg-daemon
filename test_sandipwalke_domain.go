//go:build ignore
// test_sandipwalke_domain.go
// Test EMSG setup for sandipwalke.com domain
package main

import (
	"emsg-daemon/internal/router"
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
	fmt.Println("🔍 Testing EMSG Setup for sandipwalke.com")
	fmt.Println("====================================================")

	domain := "sandipwalke.com"
	emsgDomain := "_emsg." + domain
	serverURL := "https://emsg." + domain

	// Test 1: DNS TXT Record Check
	fmt.Println("\n1. Checking DNS TXT Record...")
	fmt.Printf("Looking up: %s\n", emsgDomain)

	txtRecords, err := net.LookupTXT(emsgDomain)
	if err != nil {
		fmt.Printf("❌ DNS TXT lookup failed: %v\n", err)
		fmt.Println("   → You need to add the TXT record to your DNS")
		fmt.Printf("   → Add: %s TXT \"https://emsg.%s:8765\"\n", emsgDomain, domain)
	} else if len(txtRecords) == 0 {
		fmt.Printf("❌ No TXT records found for %s\n", emsgDomain)
		fmt.Println("   → You need to add the TXT record to your DNS")
	} else {
		fmt.Printf("✅ TXT record found: %s\n", txtRecords[0])
		for i, record := range txtRecords {
			fmt.Printf("   Record %d: %s\n", i+1, record)
		}
	}

	// Test 2: Address Validation
	fmt.Println("\n2. Testing address validation...")
	testAddress := "sandip#" + domain

	if err := router.ValidateAddress(testAddress); err != nil {
		fmt.Printf("❌ Address validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Address format valid: %s\n", testAddress)
	}

	// Test 3: Route Lookup (will work once DNS is set up)
	fmt.Println("\n3. Testing route lookup...")

	route, err := router.LookupRoute(testAddress)
	if err != nil {
		fmt.Printf("⚠️  Route lookup failed: %v\n", err)
		fmt.Println("   → This is expected if DNS TXT record is not set up yet")
	} else {
		fmt.Printf("✅ Route found: %s\n", route)

		// Parse route info
		routeInfo, err := router.ParseRouteInfo(route)
		if err != nil {
			fmt.Printf("   → Simple URL format: %s\n", route)
		} else {
			fmt.Printf("   → Server: %s\n", routeInfo.Server)
			fmt.Printf("   → Version: %s\n", routeInfo.Version)
			fmt.Printf("   → TTL: %d seconds\n", routeInfo.TTL)
		}
	}

	// Test 4: Server Connectivity (will work once server is deployed)
	fmt.Println("\n4. Testing server connectivity...")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(serverURL + "/api/user?address=test")
	if err != nil {
		fmt.Printf("⚠️  Server connection failed: %v\n", err)
		fmt.Println("   → This is expected if server is not deployed yet")
		fmt.Printf("   → Deploy EMSG daemon to: %s\n", serverURL)
	} else {
		defer resp.Body.Close()
		fmt.Printf("✅ Server responding: %d %s\n", resp.StatusCode, resp.Status)
		if resp.StatusCode == 404 {
			fmt.Println("   → EMSG daemon is running correctly!")
		}
	}

	// Test 5: HTTPS Check
	fmt.Println("\n5. Testing HTTPS configuration...")

	httpsURL := "https://emsg." + domain + "/api/user?address=test"
	resp, err = client.Get(httpsURL)
	if err != nil {
		fmt.Printf("⚠️  HTTPS connection failed: %v\n", err)
		fmt.Println("   → Set up SSL certificate with Let's Encrypt")
	} else {
		defer resp.Body.Close()
		fmt.Printf("✅ HTTPS working: %d %s\n", resp.StatusCode, resp.Status)
	}

	// Summary and Next Steps
	fmt.Println("\n====================================================")
	fmt.Println("📋 Setup Summary for sandipwalke.com")
	fmt.Println("====================================================")

	fmt.Println("\n🔧 Required Steps:")
	fmt.Println("1. Add DNS TXT record:")
	fmt.Printf("   Name: %s\n", emsgDomain)
	fmt.Println("   Type: TXT")
	fmt.Printf("   Value: https://emsg.%s:8765\n", domain)
	fmt.Println("   TTL: 3600")

	fmt.Println("\n2. Set up subdomain:")
	fmt.Printf("   Point emsg.%s to your server IP\n", domain)

	fmt.Println("\n3. Deploy EMSG daemon:")
	fmt.Println("   - Build: go build ./cmd/daemon")
	fmt.Println("   - Configure environment variables")
	fmt.Println("   - Start daemon on port 8765")

	fmt.Println("\n4. Configure SSL:")
	fmt.Println("   - Install Let's Encrypt certificate")
	fmt.Println("   - Set up Nginx reverse proxy")

	fmt.Println("\n📧 Example EMSG Addresses:")
	fmt.Printf("   - sandip#%s\n", domain)
	fmt.Printf("   - admin#%s\n", domain)
	fmt.Printf("   - contact#%s\n", domain)

	fmt.Println("\n🧪 Test Commands (after setup):")
	fmt.Printf("   dig TXT %s\n", emsgDomain)
	fmt.Printf("   curl https://emsg.%s/api/user?address=test\n", domain)

	fmt.Println("\n🎉 Once complete, your domain will be part of the EMSG network!")
}
