// boltapi.go
// REST API for EMSG Daemon using BoltDB
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"emsg-daemon/internal/auth"
	"emsg-daemon/internal/group"
	"emsg-daemon/internal/message"
	"emsg-daemon/internal/router"
	"emsg-daemon/internal/storage"

	"go.etcd.io/bbolt"
)

// BoltAPI handler struct to hold BoltDB reference
type BoltAPI struct {
	DB *bbolt.DB
}

// Example: GET /api/user?address=alice#emsg.dev
func (api *BoltAPI) ApiGetUser(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "missing address", http.StatusBadRequest)
		return
	}
	// URL decode the address to handle encoded characters like %23 for #
	decodedAddress, err := url.QueryUnescape(address)
	if err != nil {
		decodedAddress = address // fallback to original if decode fails
	}
	user, err := storage.GetUserBolt(api.DB, decodedAddress)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// Example: POST /api/user (register user with profile fields)
func (api *BoltAPI) ApiRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address        string `json:"address"`
		PubKey         string `json:"pubkey"`
		FirstName      string `json:"first_name"`
		MiddleName     string `json:"middle_name"`
		LastName       string `json:"last_name"`
		DisplayPicture string `json:"display_picture"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	user, err := auth.RegisterUser(req.Address, req.PubKey, req.FirstName, req.MiddleName, req.LastName, req.DisplayPicture)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := storage.StoreUserBolt(api.DB, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// POST /api/message (send a message)
func (api *BoltAPI) ApiSendMessage(w http.ResponseWriter, r *http.Request) {
	var msg message.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Validate message
	if err := msg.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Store message
	if err := storage.StoreMessageBolt(api.DB, &msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "message sent"})
}

// GET /api/messages?user=alice#emsg.dev (get messages for a user)
func (api *BoltAPI) ApiGetMessages(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "missing user parameter", http.StatusBadRequest)
		return
	}

	// URL decode the user address
	decodedUser, err := url.QueryUnescape(user)
	if err != nil {
		decodedUser = user
	}

	messages, err := storage.GetMessagesByUserBolt(api.DB, decodedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

// POST /api/group (create a group)
func (api *BoltAPI) ApiCreateGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		DisplayPic  string   `json:"display_pic"`
		Members     []string `json:"members"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.ID == "" || req.Name == "" {
		http.Error(w, "missing required fields: id, name", http.StatusBadRequest)
		return
	}

	grp := group.NewGroup(req.ID, req.Name, req.Description, req.DisplayPic, req.Members)

	if err := storage.StoreGroupBolt(api.DB, grp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(grp)
}

// GET /api/group?id=group1 (get group info)
func (api *BoltAPI) ApiGetGroup(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing group id", http.StatusBadRequest)
		return
	}

	grp, err := storage.GetGroupBolt(api.DB, id)
	if err != nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(grp)
}

// GET /api/route?address=alice#emsg.dev (get routing info for an address)
func (api *BoltAPI) ApiGetRoute(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "missing address parameter", http.StatusBadRequest)
		return
	}

	// URL decode the address
	decodedAddress, err := url.QueryUnescape(address)
	if err != nil {
		decodedAddress = address
	}

	// Validate address format
	if err := router.ValidateAddress(decodedAddress); err != nil {
		http.Error(w, fmt.Sprintf("invalid address: %v", err), http.StatusBadRequest)
		return
	}

	// Get route information
	routeInfo, err := router.GetRouteInfo(decodedAddress)
	if err != nil {
		http.Error(w, fmt.Sprintf("route lookup failed: %v", err), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(routeInfo)
}

// POST /api/route/validate (validate multiple addresses)
func (api *BoltAPI) ApiValidateAddresses(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Addresses []string `json:"addresses"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	results := make(map[string]interface{})

	for _, address := range req.Addresses {
		if err := router.ValidateAddress(address); err != nil {
			results[address] = map[string]interface{}{
				"valid": false,
				"error": err.Error(),
			}
		} else {
			results[address] = map[string]interface{}{
				"valid": true,
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
	})
}

// POST /api/route/message (determine routing for a message)
func (api *BoltAPI) ApiRouteMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Recipients []string `json:"recipients"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if len(req.Recipients) == 0 {
		http.Error(w, "no recipients provided", http.StatusBadRequest)
		return
	}

	routes, err := router.RouteMessage(req.Recipients)
	if err != nil {
		http.Error(w, fmt.Sprintf("routing failed: %v", err), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"routes": routes,
	})
}

// StartBoltServer starts the REST API server with BoltDB
func StartBoltServer(db *bbolt.DB, port string) {
	api := &BoltAPI{DB: db}
	auth := &AuthMiddleware{DB: db}
	// User endpoints
	http.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.ApiGetUser(w, r)
		} else if r.Method == http.MethodPost {
			api.ApiRegisterUser(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Message endpoints (protected - requires authentication)
	http.HandleFunc("/api/message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			auth.RequireAuth(api.ApiSendMessage)(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			auth.RequireAuth(api.ApiGetMessages)(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Group endpoints
	http.HandleFunc("/api/group", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.ApiGetGroup(w, r) // Public - no auth required
		} else if r.Method == http.MethodPost {
			auth.RequireAuth(api.ApiCreateGroup)(w, r) // Protected
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// DNS Routing endpoints
	http.HandleFunc("/api/route", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.ApiGetRoute(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/route/validate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api.ApiValidateAddresses(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/route/message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api.ApiRouteMessage(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	go http.ListenAndServe(":"+port, nil)
}
