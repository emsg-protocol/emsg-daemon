// storage.go
// SQLite/Postgres handling for EMSG Daemon
package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

// TODO: Add functions for storing and retrieving messages and groups

func InitDB(dataSourceName string) (*sql.DB, error) {
	return sql.Open("sqlite3", dataSourceName)
}

// StoreMessage inserts a message into the database
func StoreMessage(db *sql.DB, msg *Message) error {
	_, err := db.Exec(`INSERT INTO messages (from_addr, to_addr, cc_addr, group_id, body, signature) VALUES (?, ?, ?, ?, ?, ?)`,
		msg.From, strings.Join(msg.To, ","), strings.Join(msg.CC, ","), msg.GroupID, msg.Body, msg.Signature)
	return err
}

// GetMessagesByUser retrieves messages for a user
func GetMessagesByUser(db *sql.DB, user string) ([]Message, error) {
	rows, err := db.Query(`SELECT from_addr, to_addr, cc_addr, group_id, body, signature FROM messages WHERE to_addr LIKE ?`, "%"+user+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var msgs []Message
	for rows.Next() {
		var m Message
		var to, cc string
		if err := rows.Scan(&m.From, &to, &cc, &m.GroupID, &m.Body, &m.Signature); err != nil {
			return nil, err
		}
		m.To = strings.Split(to, ",")
		m.CC = strings.Split(cc, ",")
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// StoreGroup inserts a group into the database (with metadata and admins)
func StoreGroup(db *sql.DB, group *Group) error {
	_, err := db.Exec(`INSERT OR REPLACE INTO groups (id, name, description, display_pic, members, admins) VALUES (?, ?, ?, ?, ?, ?)`,
		group.ID, group.Name, group.Description, group.DisplayPic, strings.Join(group.Members, ","), strings.Join(group.Admins, ","))
	return err
}

// GetGroup retrieves a group by ID (with metadata and admins)
func GetGroup(db *sql.DB, id string) (*Group, error) {
	row := db.QueryRow(`SELECT id, name, description, display_pic, members, admins FROM groups WHERE id = ?`, id)
	var gid, name, desc, dp, members, admins string
	if err := row.Scan(&gid, &name, &desc, &dp, &members, &admins); err != nil {
		return nil, err
	}
	return &Group{
		ID: gid,
		Name: name,
		Description: desc,
		DisplayPic: dp,
		Members: strings.Split(members, ","),
		Admins: strings.Split(admins, ","),
	}, nil
}

// InitSchema creates the required tables if they do not exist (with group metadata and admins)
func InitSchema(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_addr TEXT, to_addr TEXT, cc_addr TEXT, group_id TEXT, body TEXT, signature TEXT
	);
	CREATE TABLE IF NOT EXISTS groups (
		id TEXT PRIMARY KEY, name TEXT, description TEXT, display_pic TEXT, members TEXT, admins TEXT
	);
	CREATE TABLE IF NOT EXISTS users (
		address TEXT PRIMARY KEY, pubkey TEXT, first_name TEXT, middle_name TEXT, last_name TEXT, display_picture TEXT
	);`)
	return err
}
