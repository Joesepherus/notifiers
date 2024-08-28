package main

import (
	"database/sql"
	"fmt"
	"sort"

	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"notifiers/controllers/alertsController"
	database "notifiers/db"
	"notifiers/services/alertsService"
	"notifiers/services/yahooService"
	"notifiers/types/alertsTypes"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

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
			alertsController.AddAlert(w, r)
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
	db = database.InitDB("./alerts.db")
	defer database.DB.Close()
	// Pass the db connection to alertsService
	alertsService.SetDB(db)

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
			untriggeredAlerts, err := alertsService.GetAlerts()

			if err != nil {
				log.Printf("Failed to fetch untriggered alerts: %v", err)
				continue
			}

			// Create a map to group alerts by symbol
			alertsBySymbol := make(map[string][]alertsTypes.Alert)

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

				stockData, err := yahooService.GetStockCurrentValue(symbol)
				if err != nil {
					log.Printf("Failed to get stock value for %s: %v", symbol, err)
					continue
				}

				currentPrice := stockData.Chart.Result[0].Meta.RegularMarketPrice
				fmt.Printf("Current price of %s: %.4f %s\n", stockData.Chart.Result[0].Meta.Symbol, currentPrice, stockData.Chart.Result[0].Meta.Currency)

				// Check and process alerts
				alertsService.CheckAlerts(symbol, currentPrice)
			}
		}
		fmt.Println("\n\n")
	}
}
