package alertsService

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func AddAlert(symbol string, triggerValue float64, alertType string) error {
	_, err := db.Exec("INSERT INTO alerts (symbol, trigger_value, alert_type) VALUES (?, ?, ?)", symbol, triggerValue, alertType)
	if err == nil {
		return nil
	}
	if err.Error() == "database is locked" {
		fmt.Printf("error inserting alerts: %v", err)
	}
	return err
}
