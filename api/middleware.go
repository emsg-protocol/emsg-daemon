// middleware.go
// Authentication middleware for EMSG Daemon API
package api

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"emsg-daemon/internal/storage"

	"go.etcd.io/bbolt"
)

// AuthMiddleware provides Ed25519 signature verification
type AuthMiddleware struct {
	DB *bbolt.DB
}

// AuthRequest represents a signed API request
type AuthRequest struct {
	Address   string `json:"address"`   // User's EMSG address
	Timestamp int64  `json:"timestamp"` // Unix timestamp
	Nonce     string `json:"nonce"`     // Random nonce to prevent replay
	Signature string `json:"signature"` // Ed25519 signature
}

// RequireAuth middleware that requires Ed25519 signature verification
func (am *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Parse Authorization header (format: "EMSG <base64-encoded-auth-request>")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "EMSG" {
			http.Error(w, "invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Decode the auth request
		authData, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			http.Error(w, "invalid authorization encoding", http.StatusUnauthorized)
			return
		}

		var authReq AuthRequest
		if err := json.Unmarshal(authData, &authReq); err != nil {
			http.Error(w, "invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Verify the signature
		if err := am.verifySignature(r, &authReq); err != nil {
			http.Error(w, fmt.Sprintf("authentication failed: %v", err), http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		r.Header.Set("X-EMSG-User", authReq.Address)

		// Call the next handler
		next(w, r)
	}
}

// verifySignature verifies the Ed25519 signature
func (am *AuthMiddleware) verifySignature(r *http.Request, authReq *AuthRequest) error {
	// Check timestamp (prevent replay attacks)
	now := time.Now().Unix()
	if authReq.Timestamp < now-300 || authReq.Timestamp > now+60 { // 5 min past, 1 min future
		return fmt.Errorf("timestamp out of range")
	}

	// Get user's public key from database
	user, err := storage.GetUserBolt(am.DB, authReq.Address)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Create the message that was signed
	signedMessage := am.createSignedMessage(r, authReq)

	// Decode the signature
	signature, err := base64.StdEncoding.DecodeString(authReq.Signature)
	if err != nil {
		return fmt.Errorf("invalid signature encoding")
	}

	// Verify the signature
	if !ed25519.Verify(user.PubKey, []byte(signedMessage), signature) {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

// createSignedMessage creates the message that should be signed
func (am *AuthMiddleware) createSignedMessage(r *http.Request, authReq *AuthRequest) string {
	// Format: METHOD:PATH:TIMESTAMP:NONCE
	return fmt.Sprintf("%s:%s:%d:%s",
		r.Method,
		r.URL.Path,
		authReq.Timestamp,
		authReq.Nonce)
}

// OptionalAuth middleware that extracts auth info if present but doesn't require it
func (am *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "EMSG" {
				authData, err := base64.StdEncoding.DecodeString(parts[1])
				if err == nil {
					var authReq AuthRequest
					if json.Unmarshal(authData, &authReq) == nil {
						if am.verifySignature(r, &authReq) == nil {
							r.Header.Set("X-EMSG-User", authReq.Address)
						}
					}
				}
			}
		}
		next(w, r)
	}
}

// CreateAuthRequest creates a signed authentication request
func CreateAuthRequest(address string, privateKey ed25519.PrivateKey, method, path string) (string, error) {
	timestamp := time.Now().Unix()
	nonce := fmt.Sprintf("%d", time.Now().UnixNano()) // Simple nonce

	// Create the message to sign
	message := fmt.Sprintf("%s:%s:%d:%s", method, path, timestamp, nonce)

	// Sign the message
	signature := ed25519.Sign(privateKey, []byte(message))

	// Create auth request
	authReq := AuthRequest{
		Address:   address,
		Timestamp: timestamp,
		Nonce:     nonce,
		Signature: base64.StdEncoding.EncodeToString(signature),
	}

	// Encode as JSON then base64
	jsonData, err := json.Marshal(authReq)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(jsonData), nil
}

// GetAuthenticatedUser extracts the authenticated user from request headers
func GetAuthenticatedUser(r *http.Request) string {
	return r.Header.Get("X-EMSG-User")
}

// IsAuthenticated checks if the request has a valid authenticated user
func IsAuthenticated(r *http.Request) bool {
	return GetAuthenticatedUser(r) != ""
}
