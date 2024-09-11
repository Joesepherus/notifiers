package loggingService

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"tradingalerts/middlewares/authMiddleware"
	"tradingalerts/utils/authUtils"
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
func LogToDB(logType, message string, r *http.Request) {
	var email string = "guest"
	var ip string
	if r != nil {
		email, _ = r.Context().Value(authMiddleware.UserEmailKey).(string)
		ip = authUtils.GetIPAddress(r)
	}

	log.Printf("User Email: %s, Type: %s, Message: %s, Endpoint: %s, IP: %s", email, logType, message, r.URL.Path, ip)
	_, err := db.Exec("INSERT INTO logs (email, type, message, endpoint, ip) VALUES ($1, $2, $3, $4, $5)", email, logType, message, r.URL.Path, ip)
	if err != nil {
		log.Printf("Error logging to DB: %v", err)
	}
	return
}
