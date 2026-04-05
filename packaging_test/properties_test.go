// properties_test.go
// Property-based tests for EMSG Daemon packaging properties.
// Uses pgregory.net/rapid for property generation.
//
// Validates: Requirements 2.3, 2.6, 3.2, 3.4, 4.1, 4.2
package packaging_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"emsg-daemon/api"
	"emsg-daemon/internal/config"
	"emsg-daemon/internal/storage"

	"pgregory.net/rapid"
)

// validLogLevels are the accepted log level values.
var validLogLevels = []string{"info", "warn", "error", "debug"}

// --- Property 2: Health Endpoint Responsiveness ---
// For any running daemon instance with a valid configuration, a GET /api/user?address=health
// request must receive an HTTP response (200 or 404) within 2 seconds.
//
// Validates: Requirements 2.3
func TestProperty2_HealthEndpointResponsiveness(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate a random valid config (port and log level only — no real server needed)
		_ = rapid.IntRange(1024, 65535).Draw(rt, "port")
		_ = rapid.SampledFrom(validLogLevels).Draw(rt, "logLevel")

		// Create a real BoltDB in a temp dir so ApiGetUser doesn't panic.
		// The health probe uses address="health" which returns 404 (user not found).
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")
		db, err := storage.InitBoltDB(dbPath)
		if err != nil {
			rt.Fatalf("InitBoltDB failed: %v", err)
		}
		defer db.Close()

		boltAPI := &api.BoltAPI{DB: db}
		mux := http.NewServeMux()
		mux.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				boltAPI.ApiGetUser(w, r)
			} else {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})

		srv := httptest.NewServer(mux)
		defer srv.Close()

		start := time.Now()
		resp, err := http.Get(srv.URL + "/api/user?address=health")
		elapsed := time.Since(start)

		if err != nil {
			rt.Fatalf("health request failed: %v", err)
		}
		defer resp.Body.Close()

		// Must respond within 2 seconds
		if elapsed > 2*time.Second {
			rt.Fatalf("health endpoint took %v, expected < 2s", elapsed)
		}

		// Must return 200 or 404
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
			rt.Fatalf("expected 200 or 404, got %d", resp.StatusCode)
		}
	})
}

// --- Property 4: Configuration Round-Trip ---
// For any valid emsg.conf file, LoadConfigFromFile must return exactly the values written.
//
// Validates: Requirements 3.2, 3.4
func TestProperty4_ConfigurationRoundTrip(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		port := rapid.IntRange(1024, 65535).Draw(rt, "port")
		logLevel := rapid.SampledFrom(validLogLevels).Draw(rt, "logLevel")
		domain := rapid.StringMatching(`[a-z]{3,10}\.[a-z]{2,5}`).Draw(rt, "domain")

		portStr := fmt.Sprintf("%d", port)

		// Write config to a temp file in KEY=VALUE format
		tmpDir := t.TempDir()
		cfgPath := filepath.Join(tmpDir, "emsg.conf")
		content := strings.Join([]string{
			"EMSG_PORT=" + portStr,
			"EMSG_LOG_LEVEL=" + logLevel,
			"EMSG_DOMAIN=" + domain,
		}, "\n") + "\n"

		if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
			rt.Fatalf("failed to write config file: %v", err)
		}

		// Load it back
		cfg, err := config.LoadConfigFromFile(cfgPath)
		if err != nil {
			rt.Fatalf("LoadConfigFromFile failed: %v", err)
		}

		// Assert all values match what was written
		if cfg.Port != portStr {
			rt.Fatalf("port round-trip failed: wrote %q, got %q", portStr, cfg.Port)
		}
		if cfg.LogLevel != logLevel {
			rt.Fatalf("log level round-trip failed: wrote %q, got %q", logLevel, cfg.LogLevel)
		}
		if cfg.Domain != domain {
			rt.Fatalf("domain round-trip failed: wrote %q, got %q", domain, cfg.Domain)
		}
	})
}

// --- Property 5: Static Asset Serving ---
// For any running daemon instance with a populated www/ directory:
//   - GET / returns index.html content with HTTP 200
//   - GET /emsg.wasm returns Content-Type: application/wasm
//   - GET /some/client/route returns index.html (SPA fallback)
//
// Validates: Requirements 4.1, 4.2
func TestProperty5_StaticAssetServing(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate random index.html content (non-empty printable ASCII)
		indexContent := rapid.StringMatching(`[A-Za-z0-9 .,!?<>/="]{10,100}`).Draw(rt, "indexContent")

		// Create a temp www/ directory with index.html and a dummy emsg.wasm
		tmpDir := t.TempDir()
		wwwDir := filepath.Join(tmpDir, "www")
		if err := os.MkdirAll(wwwDir, 0755); err != nil {
			rt.Fatalf("failed to create www dir: %v", err)
		}

		// Write index.html
		if err := os.WriteFile(filepath.Join(wwwDir, "index.html"), []byte(indexContent), 0644); err != nil {
			rt.Fatalf("failed to write index.html: %v", err)
		}

		// Write a dummy emsg.wasm (content doesn't matter for Content-Type test)
		if err := os.WriteFile(filepath.Join(wwwDir, "emsg.wasm"), []byte{0x00, 0x61, 0x73, 0x6d}, 0644); err != nil {
			rt.Fatalf("failed to write emsg.wasm: %v", err)
		}

		// Create an httptest.Server with ServeStaticAssets registered
		mux := http.NewServeMux()
		api.ServeStaticAssets(mux, wwwDir)
		srv := httptest.NewServer(mux)
		defer srv.Close()

		// Assert GET / returns index.html content with HTTP 200
		resp, err := http.Get(srv.URL + "/")
		if err != nil {
			rt.Fatalf("GET / failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			rt.Fatalf("GET / expected 200, got %d", resp.StatusCode)
		}
		buf := make([]byte, len(indexContent)+64)
		n, _ := resp.Body.Read(buf)
		body := string(buf[:n])
		if !strings.Contains(body, indexContent) {
			rt.Fatalf("GET / body does not contain index.html content: got %q", body)
		}

		// Assert GET /emsg.wasm returns Content-Type: application/wasm
		wasmResp, err := http.Get(srv.URL + "/emsg.wasm")
		if err != nil {
			rt.Fatalf("GET /emsg.wasm failed: %v", err)
		}
		defer wasmResp.Body.Close()
		ct := wasmResp.Header.Get("Content-Type")
		if ct != "application/wasm" {
			rt.Fatalf("GET /emsg.wasm Content-Type: expected application/wasm, got %q", ct)
		}

		// Assert GET /some/client/route returns index.html (SPA fallback)
		spaResp, err := http.Get(srv.URL + "/some/client/route")
		if err != nil {
			rt.Fatalf("GET /some/client/route failed: %v", err)
		}
		defer spaResp.Body.Close()
		if spaResp.StatusCode != http.StatusOK {
			rt.Fatalf("SPA fallback expected 200, got %d", spaResp.StatusCode)
		}
		spaBuf := make([]byte, len(indexContent)+64)
		spaN, _ := spaResp.Body.Read(spaBuf)
		spaBody := string(spaBuf[:spaN])
		if !strings.Contains(spaBody, indexContent) {
			rt.Fatalf("SPA fallback body does not contain index.html content: got %q", spaBody)
		}
	})
}
