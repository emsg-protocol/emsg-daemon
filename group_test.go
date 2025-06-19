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

func TestGroupMetadataAndAdmin(t *testing.T) {
	g := NewGroup("group1", "Test Group", "A test group", "http://img", []string{"alice#emsg.dev"})
	if g.Name != "Test Group" || g.Description != "A test group" || g.DisplayPic != "http://img" {
		t.Error("group metadata not set correctly")
	}
	g.AddAdmin("alice#emsg.dev")
	if len(g.Admins) != 1 || g.Admins[0] != "alice#emsg.dev" {
		t.Error("admin not added correctly")
	}
	g.UpdateName("Renamed Group")
	if g.Name != "Renamed Group" {
		t.Error("group name not updated")
	}
	g.UpdateDescription("Updated desc")
	if g.Description != "Updated desc" {
		t.Error("group description not updated")
	}
	g.UpdateDisplayPic("http://newimg")
	if g.DisplayPic != "http://newimg" {
		t.Error("group display picture not updated")
	}
	g.RemoveAdmin("alice#emsg.dev")
	if len(g.Admins) != 0 {
		t.Error("admin not removed correctly")
	}
}
