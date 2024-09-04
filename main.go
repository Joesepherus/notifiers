package main

import (
	"database/sql"
	"os"
	"sort"
	"sync"

	"log"
	"time"

	"notifiers/controllers"
	database "notifiers/db"
	"notifiers/utils/subscriptionUtils"

	// "notifiers/loadTest"
	"notifiers/payments/payments"
	"notifiers/services/alertsService"
	"notifiers/services/userService"
	"notifiers/services/yahooService"
	"notifiers/types/alertsTypes"

	"notifiers/templates"

	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	if os.Getenv("ENV") == "prod" {
		// Open or create the log file
		file, err := os.OpenFile("notifiers.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		// Set log output to the file
		log.SetOutput(file)

		// Customize the logger (optional)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.SetPrefix("INFO: ")

		// Log initialization message
		log.Printf("Initializing")
		log.Println("Logging initialized")
	}

	templates.InitTemplates()
	// start a new goroutine for the rest api endpoints
	go controllers.RestApi()
	db = database.InitDB("./alerts.db")
	defer database.DB.Close()
	// Pass the db connection to alertsService
	alertsService.SetDB(db)
	userService.SetDB(db)

	// loadTest.SetupDbWithLotsOfAlerts()
	payments.Setup()
	subscriptionUtils.Setup()

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

			// Create a wait group to wait for all goroutines to complete
			var wg sync.WaitGroup

			// Track each symbol
			for _, symbol := range symbols {
				wg.Add(1) // Increment the wait group counter
				go func(symbol string) {
					defer wg.Done() // Decrement the wait group counter when the goroutine completes

					stockData, err := yahooService.GetStockCurrentValue(symbol)
					if err != nil {
						log.Printf("Failed to get stock value for %s: %v", symbol, err)
						return
					}

					currentPrice := stockData.Chart.Result[0].Meta.RegularMarketPrice
					// fmt.Printf("Current price of %s: %.4f %s\n", stockData.Chart.Result[0].Meta.Symbol, currentPrice, stockData.Chart.Result[0].Meta.Currency)

					// Check and process alerts
					alertsService.CheckAlerts(symbol, currentPrice)
				}(symbol)
			}

			// Wait for all goroutines to finish
			wg.Wait()
		}
		// fmt.Println("\n\n")
	}
}
