// storage.go
// SQLite/Postgres handling for EMSG Daemon
package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: Add functions for storing and retrieving messages and groups

func InitDB(dataSourceName string) (*sql.DB, error) {
	return sql.Open("sqlite3", dataSourceName)
}
