package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"tradingalerts/controllers/alertsController"
	"tradingalerts/controllers/authController"
	"tradingalerts/controllers/priceChangeController"
	"tradingalerts/middlewares/authMiddleware"
	"tradingalerts/middlewares/bodySizeMiddleware"
	"tradingalerts/middlewares/logMiddleware"
	"tradingalerts/middlewares/rateLimitMiddleware"
	"tradingalerts/payments/payments"
	"tradingalerts/utils/subscriptionUtils"

	"os"
	"strconv"
	"tradingalerts/services/alertsService"
	"tradingalerts/services/loggingService"
	"tradingalerts/services/userService"
	"tradingalerts/templates"
)

// Health check handler
func healthHandler(w http.ResponseWriter) {
	// Basic health check response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	// This handler will only be called if the token is valid
	fmt.Fprintf(w, "Welcome to the protected area!")
}

func PageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET") // Specifies the HTTP methods allowed.
	w.Header().Set("X-Frame-Options", "DENY")                   // Prevents clickjacking
	w.Header().Set("X-Content-Type-Options", "nosniff")         // Prevents MIME sniffing
	w.Header().Set("X-XSS-Protection", "1; mode=block")         // Protects against XSS attacks
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
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
			log.Print("canAddAlert", subscriptionUtils.UserSubscription)
			log.Print("SubscirptionType", UserSubscription.SubscriptionType)
			log.Print("canAddAlert[email]", UserSubscription.CanAddAlert)
		}
		templateLocation = templates.BaseLocation + "/index.html"
		pageTitle = "Trading Alerts"
	case "/pricing":
		UserSubscription := subscriptionUtils.UserSubscription[email]
		data["CanAddAlert"] = UserSubscription.CanAddAlert
		data["SubscriptionType"] = UserSubscription.SubscriptionType
		log.Print("canAddAlert", subscriptionUtils.UserSubscription)
		log.Print("SubscriptionType", UserSubscription.SubscriptionType)
		log.Print("canAddAlert[email]", UserSubscription.CanAddAlert)
		templateLocation = templates.BaseLocation + "/pricing.html"
		pageTitle = "Pricing - Trading Alerts"
	case "/about":
		templateLocation = templates.BaseLocation + "/about.html"
		pageTitle = "About - Trading Alerts"
	case "/alerts":
		// Fetch alerts and add to data
		if user == nil {
			log.Println("You need to be logged in")
			loggingService.LogToDB("ERROR", "You need to be logged in", r)
			http.Redirect(w, r, "/error?message=You+need+to+be+logged+in", http.StatusSeeOther)
			return
		}
		alerts, err := alertsService.GetAlertsByUserID(user.ID)
		if err == nil {
			data["Alerts"] = alerts
		}
		completed_alerts, err := alertsService.GetCompletedAlertsByUserID(user.ID)
		if err == nil {
			data["CompletedAlerts"] = completed_alerts
		}
		UserSubscription := subscriptionUtils.UserSubscription[email]
		data["CanAddAlert"] = UserSubscription.CanAddAlert
		data["SubscriptionType"] = UserSubscription.SubscriptionType

		templateLocation = templates.BaseLocation + "/alerts.html"
		pageTitle = "Alerts - Trading Alerts"
	case "/price-change":
		templateLocation = templates.BaseLocation + "/price-change.html"
		pageTitle = "Price Change - Trading Alerts"
	case "/profile":
		if user == nil {
			log.Println("You need to be logged in")
			loggingService.LogToDB("ERROR", "You need to be logged in", r)
			http.Redirect(w, r, "/error?message=You+need+to+be+logged+in", http.StatusSeeOther)
			return
		}
		UserSubscription := subscriptionUtils.UserSubscription[email]
		data["CanAddAlert"] = UserSubscription.CanAddAlert
		data["SubscriptionType"] = UserSubscription.SubscriptionType
		data["MaxAlerts"] = subscriptionUtils.SUBSCRIPTION_LIMITS[UserSubscription.SubscriptionType]
		templateLocation = templates.BaseLocation + "/profile.html"
		pageTitle = "Profile - Trading Alerts"
	case "/reset-password-sent":
		templateLocation = templates.BaseLocation + "/reset-password-sent.html"
		pageTitle = "Reset Password - Trading Alerts"
	case "/reset-password-success":
		templateLocation = templates.BaseLocation + "/reset-password-success.html"
		pageTitle = "Reset Password Success - Trading Alerts"
	case "/subscription-success-temp":
		templateLocation = templates.BaseLocation + "/subscription-success-temp.html"
		pageTitle = "Subscription Successful - Trading Alerts"
	case "/subscription-success":
		templateLocation = templates.BaseLocation + "/subscription-success.html"
		pageTitle = "Subscription Successful - Trading Alerts"
	case "/subscription-cancel-temp":
		templateLocation = templates.BaseLocation + "/subscription-cancel-temp.html"
		pageTitle = "Subscription Cancelled - Trading Alerts"
	case "/subscription-cancel":
		templateLocation = templates.BaseLocation + "/subscription-cancel.html"
		pageTitle = "Subscription Cancelled - Trading Alerts"
	case "/token-expired":
		templateLocation = templates.BaseLocation + "/token-expired.html"
		pageTitle = "Token Expired - Trading Alerts"
	case "/docs":
		templateLocation = templates.BaseLocation + "/docs.html"
		pageTitle = "Documentation - Trading Alerts"
	case "/health":
		healthHandler(w)
		return
	case "/error":
		templateLocation = templates.BaseLocation + "/error.html"
		pageTitle = "Error - Trading Alerts"
		message := r.URL.Query().Get("message")
		data["Message"] = message
	default:
		templateLocation = templates.BaseLocation + "/404.html"
		pageTitle = "Page not found"
	}

	if ok {
		data["Email"] = email
	}
	data["Title"] = pageTitle
	data["Content"] = templateLocation

	templates.RenderTemplate(w, r, templateLocation, data)
}

func RestApi() {
	port := 8090
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	// Define the server with timeouts
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(port), // Listen on the specified port
		Handler:      nil,
		ReadTimeout:  5 * time.Second,  // Max time to read the request
		WriteTimeout: 10 * time.Second, // Max time to write the response
		IdleTimeout:  15 * time.Second, // Max time for idle connections
	}

	http.Handle("/protected", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenAuthMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(protectedHandler))))))

	http.Handle("/api/add-alert", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenAuthMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(alertsController.AddAlert))))))
	http.Handle("/api/delete-alert", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenAuthMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(alertsController.DeleteAlert))))))

	// Define the route for getting untriggered alerts
	http.Handle("/api/alerts", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenAuthMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(alertsController.GetAlerts))))))

	// Authentication routes
	http.Handle("/api/sign-up", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenCheckMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(authController.SignUp))))))
	http.Handle("/api/login", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenCheckMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(authController.Login))))))
	http.Handle("/api/logout", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenAuthMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(authController.Logout))))))
	http.Handle("/api/reset-password", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenCheckMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(authController.ResetPassword))))))
	http.Handle("/api/set-password", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenCheckMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(authController.SetPassword))))))
	http.Handle("/api/delete-account", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenCheckMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(authController.DeleteAccount))))))
	http.Handle("/api/get-hourly-change", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenCheckMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(priceChangeController.GetHourlyChange))))))

	// Stripe routes
	http.Handle("/api/create-checkout-session", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenAuthMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(payments.CreateCheckoutSession))))))
	http.Handle("/api/customer-by-email", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenAuthMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(payments.HandleGetCustomerByEmail))))))
	http.Handle("/api/cancel-subscription", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenAuthMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(payments.CancelSubscription))))))
	http.HandleFunc("/webhook", payments.HandleWebhook)
	http.Handle("/", bodySizeMiddleware.LimitRequestBodySize(authMiddleware.TokenCheckMiddleware(rateLimitMiddleware.RateLimitPerClient(logMiddleware.LogMiddleware(http.HandlerFunc(PageHandler))))))

	// Serve static files (CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/sw.js", http.FileServer(http.Dir("./")))
	http.Handle("/manifest.json", http.FileServer(http.Dir("./")))

	log.Printf("Starting server on :%d...\n", port)
	log.Fatal(server.ListenAndServe())
}
