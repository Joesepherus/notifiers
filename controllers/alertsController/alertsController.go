package alertsController

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tradingalerts/middlewares/authMiddleware"
	"tradingalerts/services/alertsService"
	"tradingalerts/services/loggingService"
	"tradingalerts/services/userService"
	"tradingalerts/services/yahooService"
	"tradingalerts/utils/errorUtils"
	"tradingalerts/utils/subscriptionUtils"
)

func AddAlert(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)
	email := r.Context().Value(authMiddleware.UserEmailKey).(string)
	user, err := userService.GetUserByEmail(email)
	if err != nil {
		log.Println("User not found")
		loggingService.LogToDB("ERROR", "User not found", r)
		http.Redirect(w, r, "/error?message=User+not+found", http.StatusSeeOther)
		return
	}

	if !subscriptionUtils.UserSubscription[email].CanAddAlert {
		log.Println("You have hit limit of active alerts")
		loggingService.LogToDB("ERROR", "You have hit limit of active alerts", r)
		http.Redirect(w, r, "/error?message=You+have+hit+limit+of+active+alerts", http.StatusSeeOther)
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
	stockData, err := yahooService.GetStockCurrentValue(yahooService.YahooBaseURL, symbol, "2m", "1d")
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

	canAddAlert, subscriptionType := subscriptionUtils.CheckToAddAlert(user.ID, email)
	subscriptionUtils.UserSubscription[email] = subscriptionUtils.UserAlertInfo{
		CanAddAlert:      canAddAlert,
		SubscriptionType: subscriptionType,
	}

	// Success response
	http.Redirect(w, r, "/alerts", http.StatusSeeOther)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response = map[string]string{"message": "Alert added successfully"}
	json.NewEncoder(w).Encode(response)
}

func GetAlerts(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)
	alerts, err := alertsService.GetAlerts()
	if err != nil {
		log.Println("Failed to fetch alerts")
		loggingService.LogToDB("ERROR", "Failed to fetch alerts", r)
		http.Redirect(w, r, "/error?message=Failed+to+fetch+alerts", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(alerts); err != nil {
		log.Println("Failed to encode alerts")
		loggingService.LogToDB("ERROR", "Failed to encode alerts", r)
		http.Redirect(w, r, "/error?message=Failed+to+encode+alerts", http.StatusSeeOther)
		return
	}
}

func DeleteAlert(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)

	email := r.Context().Value(authMiddleware.UserEmailKey).(string)
	user, err := userService.GetUserByEmail(email)
	if err != nil {
		log.Println("User not found")
		loggingService.LogToDB("ERROR", "User not found", r)
		http.Redirect(w, r, "/error?message=User+not+found", http.StatusSeeOther)
		return
	}

	// Get the ID from the query parameters
	idStr := r.URL.Query().Get("id")
	// Convert the ID to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid alert ID")
		loggingService.LogToDB("ERROR", "Invalid alert ID", r)
		http.Redirect(w, r, "/error?message=Invalid+alert+ID", http.StatusSeeOther)
		return
	}

	// Delete the alert by ID (implement your deletion logic here)
	err = alertsService.DeleteAlertByID(id)
	if err != nil {
		log.Println("Error deleting alert")
		loggingService.LogToDB("ERROR", "Error deleting alert", r)
		http.Redirect(w, r, "/error?message=Error+deleting+alert", http.StatusSeeOther)
		return
	}

	canAddAlert, subscriptionType := subscriptionUtils.CheckToAddAlert(user.ID, email)
	log.Print("canAddAlert", canAddAlert)
	subscriptionUtils.UserSubscription[email] = subscriptionUtils.UserAlertInfo{
		CanAddAlert:      canAddAlert,
		SubscriptionType: subscriptionType,
	}
	http.Redirect(w, r, "/alerts", http.StatusSeeOther)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert deleted successfully"))
}
