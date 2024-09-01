package alertsController

import (
	"encoding/json"
	"log"
	"net/http"
	"notifiers/middlewares/authMiddleware"
	"notifiers/payments/payments"
	"notifiers/services/alertsService"
	"notifiers/services/userService"
	"notifiers/services/yahooService"
	"strconv"
)

// TODO: add logic for when user is subscribed, so he can have
// more than 5 active alerts
var gold_productID string = "prod_QkzhvwCenEWmDY"
var diamond_productID string = "prod_QlltE9sAx7aY9z"

var CanAddAlert = make(map[string]bool)

const GOLD_SUBSCRIPTION_TOTAL = 100
const DIAMOND_SUBSCRIPTION_TOTAL = 1000

func CheckToAddAlert(userID int, email string) bool {
	alerts, _ := alertsService.GetAlertsByUserID(userID)

	cust, err := payments.GetCustomerByEmail(email)
	if err != nil {
		log.Printf("Error retrieving customer: %v", err)
		return false
	}
	gold_subscription, err := payments.GetSubscriptionByCustomerAndProduct(cust.ID, gold_productID)
	diamond_subscription, err2 := payments.GetSubscriptionByCustomerAndProduct(cust.ID, diamond_productID)
	log.Printf("gold_subscription", gold_subscription)
	log.Printf("diamond_subscription", diamond_subscription)
	if err == nil && gold_subscription.Status == "active" {
		if len(alerts) > GOLD_SUBSCRIPTION_TOTAL-1 {
			return false
		} else {
			return true
		}
	}

	if err2 == nil && diamond_subscription.Status == "active" {
		if len(alerts) > DIAMOND_SUBSCRIPTION_TOTAL-1 {
			return false
		} else {
			return true
		}
	}
	if len(alerts) > 4 {
		return false
	}
	return true
}

func AddAlert(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(authMiddleware.UserEmailKey).(string)
	user, err := userService.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	if !CanAddAlert[email] {
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

	canAddAlert := CheckToAddAlert(user.ID, email)
	CanAddAlert[email] = canAddAlert

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

func Setup() {
	emails := map[string]int{
		"joes@joesexperiences.com": 1,
		"test@gmail.com":           2,
	}

	for email, userID := range emails {
		CanAddAlert[email] = CheckToAddAlert(userID, email)
	}
}
