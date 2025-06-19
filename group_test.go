// group_test.go
// Tests for group management
package main

import "testing"

func TestGroupAddRemoveMember(t *testing.T) {
	g := NewGroup("group1", []string{"alice#emsg.dev"})
	if err := g.AddMember("bob#emsg.dev"); err != nil {
		t.Errorf("AddMember failed: %v", err)
	}
	if len(g.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(g.Members))
	}
	if err := g.AddMember("alice#emsg.dev"); err == nil {
		t.Error("expected error for duplicate member, got nil")
	}
	if err := g.RemoveMember("bob#emsg.dev"); err != nil {
		t.Errorf("RemoveMember failed: %v", err)
	}
	if len(g.Members) != 1 {
		t.Errorf("expected 1 member, got %d", len(g.Members))
	}
	if err := g.RemoveMember("bob#emsg.dev"); err == nil {
		t.Error("expected error for non-existent member, got nil")
	}
}
