package db

import (
	"database/sql"

	"github.com/altugbakan/card-logger/utils"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const (
	databaseFilePath = "cards.db"
	backupDirectory  = "backups"
)

func InitDB() *sql.DB {
	if db == nil {
		var err error
		db, err = sql.Open("sqlite3", databaseFilePath)
		if err != nil {
			utils.LogError("could not open the database: %v", err)
		}
	}
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func IsDBFilled() bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Sets").Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func FetchAndFillDB() {
	//TODO: Fetch from latest GitHub release
}
