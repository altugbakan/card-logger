package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/altugbakan/card-logger/utils"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func Init() *sql.DB {
	if db == nil {
		var err error
		db, err = sql.Open("sqlite", getDatabasePath())
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
	info, err := os.Stat(getDatabasePath())
	return !os.IsNotExist(err) && !(info.Size() == 0) && !info.IsDir()
}

func getDatabaseDirectory() string {
	config := utils.GetConfig()
	return filepath.Join("database", config.SetType)
}

func getDatabasePath() string {
	return filepath.Join(getDatabaseDirectory(), "cards.db")
}
