package loggingService

import (
	"database/sql"
	"log"
	"time"
)

// LogEntry stores the log information
type LogEntry struct {
	Email     string
	Type      string
	Message   string
	Endpoint  string
	Timestamp time.Time
}

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

// LogToDB saves a log entry to the SQLite database
func LogToDB(email, logType, message, endpoint, ip string) {
	log.Printf("User Email: %s, Type: %s, Message: %s, Endpoint: %s, IP: %s", email, logType, message, endpoint, ip)
	_, err := db.Exec("INSERT INTO logs (email, type, message, endpoint, ip) VALUES ($1, $2, $3, $4, $5)", email, logType, message, endpoint, ip)
	if err != nil {
		log.Printf("Error logging to DB: %v", err)
	}
	return 
}
