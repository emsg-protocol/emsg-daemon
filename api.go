// api.go
// Optional REST API for EMSG Daemon
package main

import (
	"encoding/json"
	"net/http"
)

// Example: GET /api/user?address=alice#emsg.dev
func apiGetUser(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "missing address", http.StatusBadRequest)
		return
	}
	user, err := GetUser(db, address)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// Example: POST /api/user (register user with profile fields)
func apiRegisterUser(w http.ResponseWriter, r *http.Request) {
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
	if err := StoreUser(db, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// TODO: Add more endpoints for messages, groups, etc.
