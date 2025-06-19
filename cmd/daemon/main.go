// main.go
// Starts the EMSG Daemon
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mattn/go-sqlite3"

	// Internal packages
	"emsg-daemon/internal/config"
	"emsg-daemon/internal/storage"
)

func main() {
	fmt.Println("Starting EMSG Daemon...")

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("Loaded config for domain: %s\n", cfg.Domain)

	// Initialize database
	db, err := storage.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	fmt.Println("Database initialized.")

	// Run DB schema migration/init
	if err := storage.InitSchema(db); err != nil {
		log.Fatalf("Failed to initialize DB schema: %v", err)
	}
	fmt.Println("Database schema ready.")

	// Initialize router, group, auth, message modules (stubs)
	fmt.Println("Router, group, auth, and message modules ready (stub).")

	// Start REST API server (in background)
	go func() {
		fmt.Println("Starting REST API server on :8080...")
		defer func() {
			if r := recover(); r != nil {
				log.Printf("REST API server crashed: %v", r)
			}
		}()
		if err := StartServer(db); err != nil {
			log.Printf("REST API server error: %v", err)
		}
	}()

	// Block forever (simulate daemon)
	for {
		time.Sleep(60 * time.Second)
	}
}
