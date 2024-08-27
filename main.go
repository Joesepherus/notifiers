package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"

	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"notifiers/mail"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type StockResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency           string  `json:"currency"`
				Symbol             string  `json:"symbol"`
				RegularMarketPrice float64 `json:"regularMarketPrice"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
}

var db *sql.DB

func getStockCurrentValue(symbol string) (*StockResponse, error) {
	yahooFinanceUrl := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?region=US&lang=en-US&includePrePost=false&interval=2m&useYfid=true&range=1d&corsDomain=finance.yahoo.com&.tsrc=finance", symbol)

	resp, err := http.Get(yahooFinanceUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching stock price: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var stockData StockResponse
	err = json.Unmarshal(body, &stockData)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return &stockData, nil
}

func addAlert(symbol string, triggerValue float64, alertType string) error {
	_, err := db.Exec("INSERT INTO alerts (symbol, trigger_value, alert_type) VALUES (?, ?, ?)", symbol, triggerValue, alertType)
	if err == nil {
		return nil
	}
	if err.Error() == "database is locked" {
		fmt.Printf("error inserting alerts: %v", err)
	}
	return err
}

type Alert struct {
	ID           int     `json:"id"`
	TriggerValue float64 `json:"triggerValue"`
	AlertType    string  `json:"alertType"`
	Symbol       string  `json:"symbol"`
}

func getAlertsBySymbol(symbol string) ([]Alert, error) {
	var alerts []Alert

	// Query to fetch rows from the database
	rows, err := db.Query("SELECT id, trigger_value, alert_type FROM alerts WHERE symbol = ? AND triggered = FALSE", symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %v", err)
	}
	defer rows.Close()

	// Iterate over rows and scan into struct
	for rows.Next() {
		var alert Alert
		if err := rows.Scan(&alert.ID, &alert.TriggerValue, &alert.AlertType); err != nil {
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

func getAlerts() ([]Alert, error) {
	var alerts []Alert

	// Query to fetch rows from the database
	rows, err := db.Query("SELECT id, symbol, trigger_value, alert_type FROM alerts WHERE triggered = FALSE ORDER BY symbol")
	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %v", err)
	}
	defer rows.Close()

	// Iterate over rows and scan into struct
	for rows.Next() {
		var alert Alert
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

func checkAlerts(symbol string, currentPrice float64) {
	alerts, err := getAlertsBySymbol(symbol)
	if err != nil {
		log.Fatalf("Error fetching alerts: %v", err)
	}

	for _, alert := range alerts {
		fmt.Printf("ID: %d, Trigger Value: %.4f, Alert Type: %s\n", alert.ID, alert.TriggerValue, alert.AlertType)
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
			mail.SendEmail("joes@joesexperiences.com", "Alert Triggered", fmt.Sprintf(
				"Alert triggered for %s: current price %.4f has reached the trigger value %.4f (%s)",
				symbol, currentPrice, alert.TriggerValue, alert.AlertType,
			))
		}
	}
}

func api_addAlert(w http.ResponseWriter, r *http.Request) {
	var response map[string]string

	// Decode the JSON body into an Alert struct
	var alert Alert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Invalid request body"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Example validation
	if alert.Symbol == "" || alert.TriggerValue <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Invalid alert data"}
		json.NewEncoder(w).Encode(response)
		return
	}
	stockData, err := getStockCurrentValue(alert.Symbol)
	if err != nil {
		log.Printf("Failed to get stock value for %s: %v", alert.Symbol, err)
	}
	currentPrice := stockData.Chart.Result[0].Meta.RegularMarketPrice
	var alertType string
	if currentPrice > alert.TriggerValue {
		alertType = "lower"
	} else {
		alertType = "higher"
	}
	err = addAlert(alert.Symbol, alert.TriggerValue, alertType)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Failed to store alert"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response = map[string]string{"message": "Alert added successfully"}
	json.NewEncoder(w).Encode(response)
}

func restApp() {
	port := 8089
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}
	// Serve the static HTML file
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/add-alert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api_addAlert(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	log.Printf("Starting server on :%d...\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func main() {
	var err error

	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	go restApp()
	db, err = sql.Open("sqlite3", "./alerts.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	db.SetMaxOpenConns(1)
	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS alerts (
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

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			untriggeredAlerts, err := getAlerts()
			if err != nil {
				log.Printf("Failed to fetch untriggered alerts: %v", err)
				continue
			}

			// Create a map to group alerts by symbol
			alertsBySymbol := make(map[string][]Alert)

			for _, alert := range untriggeredAlerts {
				alertsBySymbol[alert.Symbol] = append(alertsBySymbol[alert.Symbol], alert)
			}

			// Extract the symbols and sort them
			var symbols []string
			for symbol := range alertsBySymbol {
				symbols = append(symbols, symbol)
			}
			sort.Strings(symbols) // Sort symbols alphabetically

			// Track each symbol
			for _, symbol := range symbols {

				stockData, err := getStockCurrentValue(symbol)
				if err != nil {
					log.Printf("Failed to get stock value for %s: %v", symbol, err)
					continue
				}

				currentPrice := stockData.Chart.Result[0].Meta.RegularMarketPrice
				fmt.Printf("Current price of %s: %.4f %s\n", stockData.Chart.Result[0].Meta.Symbol, currentPrice, stockData.Chart.Result[0].Meta.Currency)

				// Check and process alerts
				checkAlerts(symbol, currentPrice)
			}
		}
		fmt.Println("\n\n")
	}
}
