package controllers

import (
	"fmt"
	"log"
	"net/http"
	"notifiers/controllers/alertsController"
	"notifiers/controllers/authController"
	"notifiers/middlewares/authMiddleware"
	"notifiers/payments/payments"
	"notifiers/utils/subscriptionUtils"

	"notifiers/services/alertsService"
	"notifiers/services/userService"
	"notifiers/templates"
	"os"
	"strconv"
)

func protectedHandler(w http.ResponseWriter, r *http.Request) {
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
			UserSubscription := subscriptionUtils.UserSubscription[email]
			data["CanAddAlert"] = UserSubscription.CanAddAlert
			data["SubscirptionType"] = UserSubscription.SubscriptionType
			log.Printf("canAddAlert", subscriptionUtils.UserSubscription)
			log.Printf("SubscirptionType", UserSubscription.SubscriptionType)
			log.Printf("canAddAlert[email]", UserSubscription.CanAddAlert)
		}
		templateLocation = "./templates/index.html"
		pageTitle = "Trading Alerts"
	case "/pricing":
		UserSubscription := subscriptionUtils.UserSubscription[email]
		data["CanAddAlert"] = UserSubscription.CanAddAlert
		data["SubscriptionType"] = UserSubscription.SubscriptionType
		log.Printf("canAddAlert", subscriptionUtils.UserSubscription)
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
		completed_alerts, err2 := alertsService.GetCompletedAlertsByUserID(user.ID)

		if err == nil {
			data["Alerts"] = alerts
		}
		if err2 == nil {
			data["CompletedAlerts"] = completed_alerts
		}
		UserSubscription := subscriptionUtils.UserSubscription[email]
		data["CanAddAlert"] = UserSubscription.CanAddAlert
		data["SubscriptionType"] = UserSubscription.SubscriptionType

		templateLocation = "./templates/alerts.html"
		pageTitle = "Alerts - Trading Alerts"
	case "/profile":
		UserSubscription := subscriptionUtils.UserSubscription[email]
		data["CanAddAlert"] = UserSubscription.CanAddAlert
		data["SubscriptionType"] = UserSubscription.SubscriptionType
		data["MaxAlerts"] = subscriptionUtils.SUBSCRIPTION_LIMITS[UserSubscription.SubscriptionType]
		templateLocation = "./templates/profile.html"
		pageTitle = "Profile - Trading Alerts"
	case "/reset-password-sent":
		templateLocation = "./templates/reset-password-sent.html"
		pageTitle = "Reset password - Trading Alerts"
	case "/reset-password-sucess":
		templateLocation = "./templates/reset-password-success.html"
		pageTitle = "Reset password - Trading Alerts"
	case "/subscription-success":
		templateLocation = "./templates/subscription-success.html"
		pageTitle = "Reset password - Trading Alerts"
	case "/subscription-cancel":
		templateLocation = "./templates/subscription-cancel.html"
		pageTitle = "Reset password - Trading Alerts"
	case "/token-expired":
		templateLocation = "./templates/token-expired.html"
		pageTitle = "Reset password - Trading Alerts"
	case "/docs":
		templateLocation = "./templates/docs.html"
		pageTitle = "Reset password - Trading Alerts"
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

	http.Handle("/protected", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(protectedHandler)))

	http.Handle("/api/add-alert", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(alertsController.AddAlert)))
	http.Handle("/api/delete-alert", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(alertsController.DeleteAlert)))

	// Define the route for getting untriggered alerts
	http.Handle("/api/alerts", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(alertsController.GetAlerts)))

	// Authentication routes
	http.HandleFunc("/api/sign-up", authController.SignUp)
	http.HandleFunc("/api/login", authController.Login)
	http.HandleFunc("/api/logout", authController.Logout)
	http.HandleFunc("/api/reset-password", authController.ResetPassword)
	http.HandleFunc("/api/set-password", authController.SetPassword)

	// Stripe routes
	http.HandleFunc("/api/create-checkout-session", payments.CreateCheckoutSession)
	http.HandleFunc("/api/customer-by-email", payments.HandleGetCustomerByEmail)
	http.Handle("/api/cancel-subscription", authMiddleware.TokenAuthMiddleware(http.HandlerFunc(payments.CancelSubscription)))
	http.HandleFunc("/webhook", payments.HandleWebhook)

	http.Handle("/", authMiddleware.TokenCheckMiddleware(http.HandlerFunc(pageHandler)))

	// Serve static files (CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Starting server on :%d...\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
