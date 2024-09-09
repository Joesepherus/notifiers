package db

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dataSourceName string) *sql.DB {
	var err error
	DB, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	DB.SetMaxOpenConns(1)

	statement, _ := DB.Prepare(`CREATE TABLE IF NOT EXISTS alerts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT NOT NULL,
		trigger_value REAL NOT NULL,
		alert_type TEXT CHECK(alert_type IN ('lower', 'higher')) NOT NULL,
		triggered BOOLEAN DEFAULT FALSE,
		user_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		completed_at TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`)
	statement.Exec()

	statement, _ = DB.Prepare(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	statement.Exec()

	statement, _ = DB.Prepare(`CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT,
		endpoint TEXT,
		ip TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	statement.Exec()

	if err != nil {
		log.Fatalf("Failed to create alerts table: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	return DB
}
