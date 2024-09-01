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
	"notifiers/services/userService"
	"notifiers/templates"
	"os"
	"strconv"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// This handler will only be called if the token is valid
	fmt.Fprintf(w, "Welcome to the protected area!")
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	var templateLocation, pageTitle string

	data := map[string]interface{}{
		"Email": "",
		// Add other default data here if needed
	}

	email, ok := r.Context().Value(authMiddleware.UserEmailKey).(string)
	user, err := userService.GetUserByEmail(email)

	switch r.URL.Path {
	case "/":
		data["CanAddAlert"] = false

		if err == nil {
			UserSubscription := alertsController.UserSubscription[email]
			data["CanAddAlert"] = UserSubscription.CanAddAlert
			data["SubscirptionType"] = UserSubscription.SubscriptionType
			log.Printf("canAddAlert", alertsController.UserSubscription)
			log.Printf("SubscirptionType", UserSubscription.SubscriptionType)
			log.Printf("canAddAlert[email]", UserSubscription.CanAddAlert)
		}
		templateLocation = "./templates/index.html"
		pageTitle = "Trading Alerts"
	case "/pricing":
		UserSubscription := alertsController.UserSubscription[email]
		data["CanAddAlert"] = UserSubscription.CanAddAlert
		data["SubscriptionType"] = UserSubscription.SubscriptionType
		log.Printf("canAddAlert", alertsController.UserSubscription)
		log.Printf("SubscriptionType", UserSubscription.SubscriptionType)
		log.Printf("canAddAlert[email]", UserSubscription.CanAddAlert)
		templateLocation = "./templates/pricing.html"
		pageTitle = "Pricing - Trading Alerts"
	case "/about":
		templateLocation = "./templates/about.html"
		pageTitle = "About - Trading Alerts"
	case "/alerts":
		// Fetch alerts and add to data
		alerts, err := alertsService.GetAlertsByUserID(user.ID)
		if err == nil {
			data["Alerts"] = alerts
		}
		templateLocation = "./templates/alerts.html"
		pageTitle = "Alerts - Trading Alerts"
	default:
		templateLocation = "./templates/404.html"
		pageTitle = "Page not found"
	}

	if ok {
		data["Email"] = email
	}
	data["Title"] = pageTitle
	data["Content"] = templateLocation

	templates.RenderTemplate(w, templateLocation, data)
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

	// Authentication routes
	http.HandleFunc("/api/sign-up", authController.SignUp)
	http.HandleFunc("/api/login", authController.Login)
	http.HandleFunc("/api/logout", authController.Logout)

	// Stripe routes
	http.HandleFunc("/create-checkout-session", payments.CreateCheckoutSession)
	http.HandleFunc("/webhook", payments.HandleWebhook)
	http.HandleFunc("/customer-by-email", payments.HandleGetCustomerByEmail)

	http.Handle("/", authMiddleware.TokenCheckMiddleware(http.HandlerFunc(pageHandler)))

	// Serve static files (CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Starting server on :%d...\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
