// message_test.go
// Tests for message validation and structure
package main

import "testing"

func TestMessageStruct(t *testing.T) {
	msg := Message{
		From:     "alice#emsg.dev",
		To:       []string{"bob#emsg.dev"},
		CC:       []string{"group1#emsg.chat"},
		GroupID:  "group1",
		Body:     "Hello, group!",
		Signature: "base64sig",
	}
	if msg.From == "" || len(msg.To) == 0 || msg.Body == "" {
		t.Error("required fields missing in message struct")
	}
	if msg.GroupID != "group1" {
		t.Errorf("expected group_id 'group1', got %s", msg.GroupID)
	}
}
