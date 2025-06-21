// test_env_config.go
// Test environment variable configuration
package main

import (
	"emsg-daemon/internal/config"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Testing Environment Variable Configuration...")

	// Test 1: Default configuration (no env vars)
	fmt.Println("\n1. Testing default configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Failed to load default config: %v\n", err)
		return
	}

	fmt.Printf("✅ Default config loaded:\n")
	fmt.Printf("   Domain: '%s'\n", cfg.Domain)
	fmt.Printf("   Database URL: '%s'\n", cfg.DatabaseURL)
	fmt.Printf("   Port: '%s'\n", cfg.Port)
	fmt.Printf("   Log Level: '%s'\n", cfg.LogLevel)
	fmt.Printf("   Max Connections: %d\n", cfg.MaxConnections)

	// Test 2: Set custom environment variables
	fmt.Println("\n2. Setting custom environment variables...")

	os.Setenv("EMSG_DOMAIN", "custom.emsg.dev")
	os.Setenv("EMSG_DATABASE_URL", "./custom_emsg.db")
	os.Setenv("EMSG_PORT", "9090")
	os.Setenv("EMSG_LOG_LEVEL", "debug")
	os.Setenv("EMSG_MAX_CONNECTIONS", "200")

	fmt.Println("   Set EMSG_DOMAIN=custom.emsg.dev")
	fmt.Println("   Set EMSG_DATABASE_URL=./custom_emsg.db")
	fmt.Println("   Set EMSG_PORT=9090")
	fmt.Println("   Set EMSG_LOG_LEVEL=debug")
	fmt.Println("   Set EMSG_MAX_CONNECTIONS=200")

	// Test 3: Load configuration with custom env vars
	fmt.Println("\n3. Testing custom configuration...")
	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Failed to load custom config: %v\n", err)
		return
	}

	fmt.Printf("✅ Custom config loaded:\n")
	fmt.Printf("   Domain: '%s'\n", cfg.Domain)
	fmt.Printf("   Database URL: '%s'\n", cfg.DatabaseURL)
	fmt.Printf("   Port: '%s'\n", cfg.Port)
	fmt.Printf("   Log Level: '%s'\n", cfg.LogLevel)
	fmt.Printf("   Max Connections: %d\n", cfg.MaxConnections)

	// Verify the values are correct
	if cfg.Domain == "custom.emsg.dev" {
		fmt.Println("   ✅ Domain correctly set from environment")
	} else {
		fmt.Printf("   ❌ Domain incorrect: expected 'custom.emsg.dev', got '%s'\n", cfg.Domain)
	}

	if cfg.DatabaseURL == "./custom_emsg.db" {
		fmt.Println("   ✅ Database URL correctly set from environment")
	} else {
		fmt.Printf("   ❌ Database URL incorrect: expected './custom_emsg.db', got '%s'\n", cfg.DatabaseURL)
	}

	if cfg.Port == "9090" {
		fmt.Println("   ✅ Port correctly set from environment")
	} else {
		fmt.Printf("   ❌ Port incorrect: expected '9090', got '%s'\n", cfg.Port)
	}

	// Test 4: Test partial environment variables
	fmt.Println("\n4. Testing partial environment variables...")

	os.Unsetenv("EMSG_PORT") // Remove port, should use default
	os.Setenv("EMSG_DOMAIN", "partial.emsg.dev")

	fmt.Println("   Unset EMSG_PORT (should use default)")
	fmt.Println("   Set EMSG_DOMAIN=partial.emsg.dev")

	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Failed to load partial config: %v\n", err)
		return
	}

	fmt.Printf("✅ Partial config loaded:\n")
	fmt.Printf("   Domain: '%s'\n", cfg.Domain)
	fmt.Printf("   Database URL: '%s'\n", cfg.DatabaseURL)
	fmt.Printf("   Port: '%s'\n", cfg.Port)

	// Test 5: Clear all environment variables
	fmt.Println("\n5. Testing with cleared environment variables...")

	os.Unsetenv("EMSG_DOMAIN")
	os.Unsetenv("EMSG_DATABASE_URL")
	os.Unsetenv("EMSG_PORT")

	fmt.Println("   Cleared all EMSG environment variables")

	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Failed to load cleared config: %v\n", err)
		return
	}

	fmt.Printf("✅ Cleared config loaded (should be defaults):\n")
	fmt.Printf("   Domain: '%s'\n", cfg.Domain)
	fmt.Printf("   Database URL: '%s'\n", cfg.DatabaseURL)
	fmt.Printf("   Port: '%s'\n", cfg.Port)

	fmt.Println("\n🎉 Environment variable configuration testing completed!")

	// Set up recommended configuration for production
	fmt.Println("\n📋 Recommended environment variables for production:")
	fmt.Println("   export EMSG_DOMAIN=yourdomain.com")
	fmt.Println("   export EMSG_DATABASE_URL=./emsg_production.db")
	fmt.Println("   export EMSG_PORT=8080")
	fmt.Println("\nOr for Windows PowerShell:")
	fmt.Println("   $env:EMSG_DOMAIN=\"yourdomain.com\"")
	fmt.Println("   $env:EMSG_DATABASE_URL=\"./emsg_production.db\"")
	fmt.Println("   $env:EMSG_PORT=\"8080\"")
}
