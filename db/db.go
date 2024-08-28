package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dataSourceName string) *sql.DB {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	DB.SetMaxOpenConns(1)
	return DB
}
