// group.go
// Persistent group system for EMSG Daemon
package main

import "errors"

type Group struct {
	ID      string
	Members []string // user addresses
}

// NewGroup creates a new group with given ID and members
func NewGroup(id string, members []string) *Group {
	return &Group{ID: id, Members: members}
}

// AddMember adds a user to the group
func (g *Group) AddMember(address string) error {
	for _, m := range g.Members {
		if m == address {
			return errors.New("user already in group")
		}
	}
	g.Members = append(g.Members, address)
	return nil
}

// RemoveMember removes a user from the group
func (g *Group) RemoveMember(address string) error {
	for i, m := range g.Members {
		if m == address {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found in group")
}
