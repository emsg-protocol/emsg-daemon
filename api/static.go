// static.go
// Static file serving for the EMSG Daemon HTTP server.
package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ServeStaticAssets registers static file routes on the given mux.
// wwwDir is the absolute path to the UI dist/ directory.
// If wwwDir does not exist, logs a warning and skips registration (API routes unaffected).
func ServeStaticAssets(mux *http.ServeMux, wwwDir string) {
	if _, err := os.Stat(wwwDir); os.IsNotExist(err) {
		log.Printf("WARNING: www directory not found at %s — UI will not be served", wwwDir)
		return
	}

	// Serve emsg.wasm with correct Content-Type
	mux.HandleFunc("/emsg.wasm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/wasm")
		http.ServeFile(w, r, filepath.Join(wwwDir, "emsg.wasm"))
	})

	// Serve wasm_exec.js
	mux.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(wwwDir, "wasm_exec.js"))
	})

	// Serve /assets/ via file server
	mux.Handle("/assets/", http.FileServer(http.Dir(wwwDir)))

	// SPA fallback: catch-all for "/" — serves index.html for any non-API path
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(wwwDir, "index.html"))
	})
}
