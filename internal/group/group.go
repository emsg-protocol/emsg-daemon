// group.go
// Persistent group system for EMSG Daemon
package main

import (
	"database/sql"
	"errors"
)

type Group struct {
	ID          string
	Members     []string // user addresses
	Admins      []string // admin addresses
	Name        string // group_name
	Description string // group_description
	DisplayPic  string // group_display_picture (URL or hash)
}

// NewGroup creates a new group with given metadata and members
func NewGroup(id, name, description, displayPic string, members []string) *Group {
	return &Group{ID: id, Name: name, Description: description, DisplayPic: displayPic, Members: members}
}

// AddAdmin assigns admin rights to a user and triggers a system message
func (g *Group) AddAdmin(address string) {
	for _, a := range g.Admins {
		if a == address {
			return // already admin
		}
	}
	g.Admins = append(g.Admins, address)
	SendSystemMessage(SystemAdminAssigned, g.ID, address)
}

// RemoveAdmin revokes admin rights from a user and triggers a system message
func (g *Group) RemoveAdmin(address string) {
	for i, a := range g.Admins {
		if a == address {
			g.Admins = append(g.Admins[:i], g.Admins[i+1:]...)
			SendSystemMessage(SystemAdminRevoked, g.ID, address)
			return
		}
	}
}

// AddMember adds a user to the group and triggers a system message
func (g *Group) AddMember(address string) error {
	for _, m := range g.Members {
		if m == address {
			return errors.New("user already in group")
		}
	}
	g.Members = append(g.Members, address)
	SendSystemMessage(SystemUserJoined, g.ID, address)
	return nil
}

// RemoveMember removes a user from the group and triggers a system message
func (g *Group) RemoveMember(address string) error {
	for i, m := range g.Members {
		if m == address {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			SendSystemMessage(SystemUserLeft, g.ID, address)
			return nil
		}
	}
	return errors.New("user not found in group")
}

// RemoveUserByAdmin removes a user by admin action and triggers a system message
func (g *Group) RemoveUserByAdmin(address string) error {
	for i, m := range g.Members {
		if m == address {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			SendSystemMessage(SystemUserRemoved, g.ID, address)
			return nil
		}
	}
	return errors.New("user not found in group")
}

// UpdateName updates the group's name and triggers a system message
func (g *Group) UpdateName(newName string) {
	g.Name = newName
	SendSystemMessage(SystemGroupRenamed, g.ID, "")
}

// UpdateDescription updates the group's description and triggers a system message
func (g *Group) UpdateDescription(newDesc string) {
	g.Description = newDesc
	SendSystemMessage(SystemDescriptionUpdated, g.ID, "")
}

// UpdateDisplayPic updates the group's display picture and triggers a system message
func (g *Group) UpdateDisplayPic(newDP string) {
	g.DisplayPic = newDP
	SendSystemMessage(SystemDPUpdated, g.ID, "")
}

// Persist group state using storage.go
func (g *Group) Save(db *sql.DB) error {
	return StoreGroup(db, g)
}

// LoadGroup loads a group from storage
func LoadGroup(db *sql.DB, id string) (*Group, error) {
	return GetGroup(db, id)
}
