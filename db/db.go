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

	statement, _ := DB.Prepare(`CREATE TABLE IF NOT EXISTS alerts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT NOT NULL,
		trigger_value REAL NOT NULL,
		alert_type TEXT CHECK(alert_type IN ('lower', 'higher')) NOT NULL,
		triggered BOOLEAN DEFAULT FALSE
	);`)
	statement.Exec()
	if err != nil {
		log.Fatalf("Failed to create alerts table: %v", err)
	}

	return DB
}
