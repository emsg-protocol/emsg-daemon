// message.go
// Core message handling for EMSG Daemon
package message

import (
	"encoding/base64"
	"errors"

	"emsg-daemon/internal/group"
	"emsg-daemon/internal/auth"
)

type Message struct {
	From      string   `json:"from"`
	To        []string `json:"to"`
	CC        []string `json:"cc"`
	GroupID   string   `json:"group_id"`
	Body      string   `json:"body"`
	Signature string   `json:"signature"`
}

// Validate checks if the message has required fields
func (m *Message) Validate() error {
	if m.From == "" || len(m.To) == 0 || m.Body == "" {
		return errors.New("missing required fields: from, to, or body")
	}
	return nil
}

// Verify checks the message signature using the sender's public key
func (m *Message) Verify(pubKey []byte) bool {
	// Use VerifySignature from auth.go
	return auth.VerifySignature(pubKey, []byte(m.Body), decodeBase64(m.Signature))
}

// Deliver delivers the message to all recipients and group members
func (m *Message) Deliver(groups map[string]*group.Group) ([]string, error) {
	delivered := make(map[string]struct{})
	for _, to := range m.To {
		delivered[to] = struct{}{}
	}
	for _, cc := range m.CC {
		if group, ok := groups[cc]; ok {
			for _, member := range group.Members {
				delivered[member] = struct{}{}
			}
		} else {
			delivered[cc] = struct{}{}
		}
	}
	if m.GroupID != "" {
		if group, ok := groups[m.GroupID]; ok {
			for _, member := range group.Members {
				delivered[member] = struct{}{}
			}
		}
	}
	var recipients []string
	for addr := range delivered {
		recipients = append(recipients, addr)
	}
	return recipients, nil
}

func decodeBase64(s string) []byte {
	// Helper for base64 decoding
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}

// TODO: Add message delivery logic
