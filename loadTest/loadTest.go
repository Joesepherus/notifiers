package loadTest

import (
	"time"
	"tradingalerts/services/alertsService"

	"golang.org/x/exp/rand"
)

func SetupDbWithLotsOfAlerts() {
	// Seed the random number generator
	rand.Seed(uint64(time.Now().UnixNano()))

	// Define the list of possible symbols
	symbols := []string{"MXN=X", "USDJPY=X", "NVDA", "AAPL", "AI", "META"}

	// Loop 10,000 times
	for i := 0; i < 100000; i++ {
		// Generate a random number between 30 and 200
		randomValue := float64(rand.Intn(4000) + 3000)

		// Select a random symbol from the list
		randomSymbol := symbols[rand.Intn(len(symbols))]

		// Call AddAlert with the random value
		alertsService.AddAlert(1, randomSymbol, randomValue, "higher")
	}
}
