package main

import (
	"database/sql"
	"fmt"
	"sort"

	"log"
	"time"

	"notifiers/controllers"
	database "notifiers/db"

	// "notifiers/loadTest"
	"notifiers/services/alertsService"
	"notifiers/services/userService"
	"notifiers/services/yahooService"
	"notifiers/types/alertsTypes"

	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	var err error = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// start a new goroutine for the rest api endpoints
	go controllers.RestApi()
	db = database.InitDB("./alerts.db")
	defer database.DB.Close()
	// Pass the db connection to alertsService
	alertsService.SetDB(db)
	userService.SetDB(db)

	// loadTest.SetupDbWithLotsOfAlerts()

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
				go alertsService.CheckAlerts(symbol, currentPrice)
			}
		}
		fmt.Println("\n\n")
	}
}
