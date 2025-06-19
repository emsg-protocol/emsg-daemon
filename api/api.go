// api.go
// Optional REST API for EMSG Daemon
package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// API handler struct to hold DB reference
type API struct {
	db *sql.DB
}

// Example: GET /api/user?address=alice#emsg.dev
func (api *API) apiGetUser(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "missing address", http.StatusBadRequest)
		return
	}
	user, err := GetUser(api.db, address)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// Example: POST /api/user (register user with profile fields)
func (api *API) apiRegisterUser(w http.ResponseWriter, r *http.Request) {
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
	user, err := RegisterUser(req.Address, req.PubKey, req.FirstName, req.MiddleName, req.LastName, req.DisplayPicture)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := StoreUser(api.db, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// StartServer starts the REST API server
func StartServer(db *sql.DB) {
	api := &API{db: db}
	http.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.apiGetUser(w, r)
		} else if r.Method == http.MethodPost {
			api.apiRegisterUser(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// TODO: Add more endpoints for messages, groups, etc.
	go http.ListenAndServe(":8080", nil)
}
