// system.go
// Handles system-generated messages for EMSG Daemon
package main

import (
	"database/sql"
	"fmt"
	"time"
)

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

// SendSystemMessage creates, stores, and logs a system message for group events
func SendSystemMessage(event, groupID, user string, db *sql.DB) {
	msg := fmt.Sprintf("[SYSTEM] %s: user %s in group %s", event, user, groupID)
	fmt.Println(msg)
	if db != nil {
		db.Exec(`INSERT INTO messages (from_addr, to_addr, cc_addr, group_id, body, signature) VALUES (?, ?, ?, ?, ?, ?)`,
			"system#local", groupID, "", groupID, msg, "")
	}
	// TODO: Broadcast to group members if needed
}
