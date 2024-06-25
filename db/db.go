package db

import (
	"database/sql"
	"os"

	"github.com/altugbakan/card-logger/utils"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const (
	backupDirectory = "backups"
)

func Init() *sql.DB {
	if db == nil {
		var err error
		db, err = sql.Open("sqlite3", utils.DatabaseFilePath)
		if err != nil {
			utils.LogError("could not open the database: %v", err)
		}
	}
	return db
}

func Close() {
	if db != nil {
		db.Close()
		db = nil
	}
}

func Reinit() {
	Close()
	db = Init()
}

func Exists() bool {
	info, err := os.Stat(utils.DatabaseFilePath)
	return !os.IsNotExist(err) && !(info.Size() == 0) && !info.IsDir()
}
