// main.go
// Starts the EMSG Daemon with BoltDB (pure Go, no CGO required)
package main

import (
	"fmt"
	"log"
	"time"

	// Internal packages
	"emsg-daemon/api"
	"emsg-daemon/internal/config"
	"emsg-daemon/internal/storage"
)

func main() {
	fmt.Println("Starting EMSG Daemon with BoltDB...")

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Loaded config for domain: %s\n", cfg.Domain)

	// Set default database URL if not provided
	dbURL := cfg.DatabaseURL
	if dbURL == "" {
		dbURL = "./emsg.db"
		fmt.Printf("Using default database: %s\n", dbURL)
	}

	// Initialize BoltDB database (pure Go, no CGO required)
	db, err := storage.InitBoltDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize BoltDB: %v", err)
	}
	defer db.Close()
	fmt.Println("BoltDB database initialized.")

	// Initialize router, group, auth, message modules (stubs)
	fmt.Println("Router, group, auth, and message modules ready (stub).")

	// Start REST API server (in background)
	go func() {
		fmt.Printf("Starting REST API server on :%s...\n", cfg.Port)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("REST API server crashed: %v", r)
			}
		}()
		api.StartBoltServer(db, cfg.Port)
	}()

	fmt.Println("EMSG Daemon is running. Press Ctrl+C to stop.")

	// Block forever (simulate daemon)
	for {
		time.Sleep(60 * time.Second)
	}
}

// Docker build command
// RUN docker build -t emsg-daemon .
// Docker run command
// docker run -p 8080:8080 emsg-daemon
