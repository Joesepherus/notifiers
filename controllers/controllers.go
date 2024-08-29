package controllers

import (
	"log"
	"net/http"
	"notifiers/controllers/alertsController"
	"notifiers/services/alertsService"
	"os"
	"strconv"
	"text/template"
)

func RestApi() {
	port := 8089
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	// Serve the static HTML file
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Define the route and its handler
	http.HandleFunc("/api/add-alert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			alertsController.AddAlert(w, r) // Assuming AddAlert is in the same package
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Define the route for getting untriggered alerts
	http.HandleFunc("/api/alerts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			alertsController.GetAlerts(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Define the route for rendering the alerts page
	http.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetAlertsPage(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Starting server on :%d...\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func GetAlertsPage(w http.ResponseWriter, r *http.Request) {
	alerts, err := alertsService.GetAlerts()
	if err != nil {
		http.Error(w, "Failed to fetch alerts", http.StatusInternalServerError)
	}

	tmpl, err := template.ParseFiles("templates/alerts.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, alerts); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
