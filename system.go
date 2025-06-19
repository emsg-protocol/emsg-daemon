// system.go
// Handles system-generated messages for EMSG Daemon
package main

import "fmt"

// Expanded system message types
const (
	SystemGroupCreated       = "group_created"
	SystemUserJoined         = "user_joined"
	SystemUserLeft           = "user_left"
	SystemUserRemoved        = "user_removed"
	SystemAdminAssigned      = "admin_assigned"
	SystemAdminRevoked       = "admin_revoked"
	SystemGroupRenamed       = "group_renamed"
	SystemDescriptionUpdated = "description_updated"
	SystemDPUpdated          = "dp_updated"
)

// SendSystemMessage creates and logs a system message for group events
func SendSystemMessage(event, groupID, user string) {
	msg := fmt.Sprintf("[SYSTEM] %s: user %s in group %s", event, user, groupID)
	fmt.Println(msg)
	// Optionally, store system messages in the database or broadcast to group members
}
