package alertsController

import (
	"encoding/json"
	"log"
	"net/http"
	"notifiers/services/alertsService"
	"notifiers/services/yahooService"
	"notifiers/types/alertsTypes"
)

func AddAlert(w http.ResponseWriter, r *http.Request) {
	var response map[string]string

	// Decode the JSON body into an Alert struct
	var alert alertsTypes.Alert
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
	stockData, err := yahooService.GetStockCurrentValue(alert.Symbol)
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
	err = alertsService.AddAlert(alert.Symbol, alert.TriggerValue, alertType)
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
