package loggingService

import (
	"database/sql"
	"log"
	"time"
)

// LogEntry stores the log information
type LogEntry struct {
	Email     string
	Endpoint  string
	Timestamp time.Time
}

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

// LogToDB saves a log entry to the SQLite database
func LogToDB(email, endpoint string, ip string) error {
	log.Printf("User Email: %s, Endpoint: %s, IP: %s", email, endpoint, ip)
	_, err := db.Exec("INSERT INTO logs (email, endpoint, ip) VALUES (?, ?, ?)", email, endpoint, ip)
	if err != nil {
		log.Printf("Error logging to DB: %v", err)
		return err
	}
	return nil
}
