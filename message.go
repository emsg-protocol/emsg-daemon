// message.go
// Core message handling for EMSG Daemon
package main

type Message struct {
	From     string   `json:"from"`
	To       []string `json:"to"`
	CC       []string `json:"cc"`
	GroupID  string   `json:"group_id"`
	Body     string   `json:"body"`
	Signature string  `json:"signature"`
}

// TODO: Add message validation, signature checking, and delivery logic
