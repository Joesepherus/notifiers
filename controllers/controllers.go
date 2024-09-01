package controllers

import (
	"fmt"
	"log"
	"net/http"
	"notifiers/controllers/alertsController"
	"notifiers/controllers/authController"
	"notifiers/middlewares/authMiddleware"
	"notifiers/payments/payments"
	"notifiers/services/alertsService"
	"notifiers/templates"
	"os"
	"strconv"
	"text/template"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// This handler will only be called if the token is valid
	fmt.Fprintf(w, "Welcome to the protected area!")
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	var templateLocation, pageTitle string

	switch r.URL.Path {
	case "/":
		templateLocation = "./templates/index.html"
		pageTitle = "Trading Alerts"
	case "/pricing":
		templateLocation = "./templates/pricing.html"
		pageTitle = "Pricing - Trading Alerts"
	case "/about":
		templateLocation = "./templates/about.html"
		pageTitle = "About - Trading Alerts"
	default:
		templateLocation = "./templates/404.html"
		pageTitle = "Page not found"
	}

	email := r.Context().Value(authMiddleware.UserEmailKey).(string)

	templates.RenderTemplate(w, templateLocation, pageTitle, email)
}

func RestApi() {
	port := 8089
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	http.Handle("/protected", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(ProtectedHandler)))

	http.Handle("/api/add-alert", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(alertsController.AddAlert)))

	// Define the route for getting untriggered alerts
	http.Handle("/api/alerts", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(alertsController.GetAlerts)))

	http.Handle("/alerts", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(GetAlertsPage)))

	// Define the route for serving the signup page
	http.HandleFunc("/sign-up", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "templates/signup.html")
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Define the route for serving the login page
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "templates/login.html")
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Authentication routes
	http.HandleFunc("/api/sign-up", authController.SignUp)
	http.HandleFunc("/api/login", authController.Login)

	// Stripe routes
	http.HandleFunc("/create-checkout-session", payments.CreateCheckoutSession)
	http.HandleFunc("/webhook", payments.HandleWebhook)

	http.Handle("/", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(pageHandler)))

	// Serve static files (CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Starting server on :%d...\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func GetAlertsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
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
