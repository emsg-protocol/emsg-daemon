// boltdb.go
// BoltDB implementation for EMSG Daemon (pure Go, no CGO required)
package storage

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.etcd.io/bbolt"
	"emsg-daemon/internal/message"
	"emsg-daemon/internal/group"
	"emsg-daemon/internal/auth"
)

var (
	messagesBucket = []byte("messages")
	groupsBucket   = []byte("groups")
	usersBucket    = []byte("users")
)

// InitBoltDB initializes a BoltDB database
func InitBoltDB(dataSourceName string) (*bbolt.DB, error) {
	db, err := bbolt.Open(dataSourceName, 0600, nil)
	if err != nil {
		return nil, err
	}

	// Create buckets if they don't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(messagesBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(groupsBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(usersBucket); err != nil {
			return err
		}
		return nil
	})

	return db, err
}

// StoreMessageBolt stores a message in BoltDB
func StoreMessageBolt(db *bbolt.DB, msg *message.Message) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(messagesBucket)
		
		// Generate a simple key (could be improved with proper ID generation)
		key := fmt.Sprintf("%s_%s_%d", msg.From, strings.Join(msg.To, ","), len(msg.Body))
		
		data, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		
		return b.Put([]byte(key), data)
	})
}

// GetMessagesByUserBolt retrieves messages for a user from BoltDB
func GetMessagesByUserBolt(db *bbolt.DB, user string) ([]message.Message, error) {
	var messages []message.Message
	
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(messagesBucket)
		
		return b.ForEach(func(k, v []byte) error {
			var msg message.Message
			if err := json.Unmarshal(v, &msg); err != nil {
				return err
			}
			
			// Check if user is in To or CC fields
			for _, to := range msg.To {
				if strings.Contains(to, user) {
					messages = append(messages, msg)
					return nil
				}
			}
			for _, cc := range msg.CC {
				if strings.Contains(cc, user) {
					messages = append(messages, msg)
					return nil
				}
			}
			
			return nil
		})
	})
	
	return messages, err
}

// StoreGroupBolt stores a group in BoltDB
func StoreGroupBolt(db *bbolt.DB, grp *group.Group) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(groupsBucket)
		
		data, err := json.Marshal(grp)
		if err != nil {
			return err
		}
		
		return b.Put([]byte(grp.ID), data)
	})
}

// GetGroupBolt retrieves a group from BoltDB
func GetGroupBolt(db *bbolt.DB, id string) (*group.Group, error) {
	var grp group.Group
	
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(groupsBucket)
		
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("group not found: %s", id)
		}
		
		return json.Unmarshal(data, &grp)
	})
	
	if err != nil {
		return nil, err
	}
	
	return &grp, nil
}

// StoreUserBolt stores a user in BoltDB
func StoreUserBolt(db *bbolt.DB, user *auth.User) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(usersBucket)
		
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		
		return b.Put([]byte(user.Address), data)
	})
}

// GetUserBolt retrieves a user from BoltDB
func GetUserBolt(db *bbolt.DB, address string) (*auth.User, error) {
	var user auth.User
	
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(usersBucket)
		
		data := b.Get([]byte(address))
		if data == nil {
			return fmt.Errorf("user not found: %s", address)
		}
		
		return json.Unmarshal(data, &user)
	})
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}
