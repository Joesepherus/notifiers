package alertsService

import (
	"database/sql"
	"fmt"
	"log"
	"notifiers/mail"
	"notifiers/types/alertsTypes"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func AddAlert(userID int, symbol string, triggerValue float64, alertType string) error {
	result, err := db.Exec("INSERT INTO alerts (user_id, symbol, trigger_value, alert_type) VALUES (?, ?, ?, ?)", userID, symbol, triggerValue, alertType)
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Alert created with id: %d", lastInsertID)
	if err == nil {
		return nil
	}
	if err.Error() == "database is locked" {
		fmt.Printf("error inserting alerts: %v", err)
	}
	return err
}

func GetAlerts() ([]alertsTypes.Alert, error) {
	var alerts []alertsTypes.Alert

	// Query to fetch rows from the database
	rows, err := db.Query("SELECT id, symbol, trigger_value, alert_type FROM alerts WHERE triggered = FALSE ORDER BY symbol")

	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %v", err)
	}
	defer rows.Close()

	// Iterate over rows and scan into struct
	for rows.Next() {
		var alert alertsTypes.Alert
		if err := rows.Scan(&alert.ID, &alert.Symbol, &alert.TriggerValue, &alert.AlertType); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		alerts = append(alerts, alert)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return alerts, nil
}

func GetAlertsByUserID(userID int) ([]alertsTypes.Alert, error) {
	var alerts []alertsTypes.Alert

	// Query to fetch rows from the database
	rows, err := db.Query(`SELECT id, symbol, trigger_value, alert_type 
		FROM alerts 
		WHERE triggered = FALSE AND user_id = $1 
		ORDER BY symbol`, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %v", err)
	}
	defer rows.Close()

	// Iterate over rows and scan into struct
	for rows.Next() {
		var alert alertsTypes.Alert
		if err := rows.Scan(&alert.ID, &alert.Symbol, &alert.TriggerValue, &alert.AlertType); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		alerts = append(alerts, alert)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return alerts, nil
}

func CheckAlerts(symbol string, currentPrice float64) {
	alerts, err := GetAlertsBySymbol(symbol)
	if err != nil {
		log.Fatalf("Error fetching alerts: %v", err)
	}

	for _, alert := range alerts {
		// fmt.Printf("ID: %d, Symbol: %s, Trigger Value: %.4f, Alert Type: %s\n", alert.ID, alert.Symbol, alert.TriggerValue, alert.AlertType)
		var shouldTrigger bool
		if alert.AlertType == "higher" && currentPrice >= alert.TriggerValue {
			shouldTrigger = true
		} else if alert.AlertType == "lower" && currentPrice <= alert.TriggerValue {
			shouldTrigger = true
		}

		if shouldTrigger {
			log.Printf("Alert triggered for %s: current price %.4f has reached the trigger value %.4f (%s)", symbol, currentPrice, alert.TriggerValue, alert.AlertType)

			statement, err := db.Prepare("UPDATE alerts SET triggered = TRUE WHERE id = ?")
			fmt.Printf("id: ", alert.ID)

			if err != nil {
				fmt.Printf("Wtf error\n", err)
				return
			}
			_, err = statement.Exec(alert.ID)
			if err != nil {
				fmt.Printf("Hey error\n", err)
				return
			} else {

				fmt.Printf("Updated\n")
			}
			go mail.SendEmail("joes@joesexperiences.com", "Alert Triggered", fmt.Sprintf(
				"Alert triggered for %s: current price %.4f has reached the trigger value %.4f (%s)",
				symbol, currentPrice, alert.TriggerValue, alert.AlertType,
			))
		}
	}
}

func GetAlertsBySymbol(symbol string) ([]alertsTypes.Alert, error) {
	var alerts []alertsTypes.Alert

	// Query to fetch rows from the database
	rows, err := db.Query("SELECT id, symbol, trigger_value, alert_type FROM alerts WHERE symbol = ? AND triggered = FALSE", symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %v", err)
	}
	defer rows.Close()

	// Iterate over rows and scan into struct
	for rows.Next() {
		var alert alertsTypes.Alert
		if err := rows.Scan(&alert.ID, &alert.Symbol, &alert.TriggerValue, &alert.AlertType); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		alerts = append(alerts, alert)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return alerts, nil
}
