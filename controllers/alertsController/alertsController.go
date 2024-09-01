package alertsController

import (
	"encoding/json"
	"log"
	"net/http"
	"notifiers/middlewares/authMiddleware"
	"notifiers/services/alertsService"
	"notifiers/services/userService"
	"notifiers/services/yahooService"
	"strconv"
)

// TODO: add logic for when user is subscribed, so he can have
// more than 5 active alerts
func AddAlert(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(authMiddleware.UserEmailKey).(string)
	user, err := userService.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	alerts, _ := alertsService.GetAlertsByUserID(user.ID)
	if len(alerts) > 4 {
		http.Error(w, "You have hit limit of 5 active alerts for free tier.", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var response map[string]string

	// Parse form values
	err = r.ParseForm()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Failed to parse form data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Extract form values
	symbol := r.FormValue("symbol")
	triggerValueStr := r.FormValue("triggerValue")

	// Convert triggerValue to float64
	triggerValue, err := strconv.ParseFloat(triggerValueStr, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Invalid trigger value"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Example validation
	if symbol == "" || triggerValue <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response = map[string]string{"error": "Invalid alert data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Fetch stock data
	stockData, err := yahooService.GetStockCurrentValue(symbol)
	if err != nil {
		log.Printf("Failed to get stock value for %s: %v", symbol, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response = map[string]string{"error": "Failed to get stock data"}
		json.NewEncoder(w).Encode(response)
		return
	}

	currentPrice := stockData.Chart.Result[0].Meta.RegularMarketPrice
	var alertType string
	if currentPrice > triggerValue {
		alertType = "lower"
	} else {
		alertType = "higher"
	}

	// Add alert to the database
	err = alertsService.AddAlert(user.ID, symbol, triggerValue, alertType)
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

func GetAlerts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	alerts, err := alertsService.GetAlerts()
	if err != nil {
		http.Error(w, "Failed to fetch alerts", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(alerts); err != nil {
		http.Error(w, "Failed to encode alerts", http.StatusInternalServerError)
	}
}
