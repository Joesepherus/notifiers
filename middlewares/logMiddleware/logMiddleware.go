package logMiddleware

import (
	"log"
	"net/http"
	"strings"
	"tradingalerts/middlewares/authMiddleware"
	"tradingalerts/services/loggingService"
)

// Helper function to get the client's real IP address, including proxies
func getIPAddress(r *http.Request) string {
	ip := r.RemoteAddr
	// Check if the request is coming from a proxy and get the real client IP
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		ip = strings.Split(forwardedFor, ",")[0] // Get the first IP in the chain
	}
	return ip
}

// LogMiddleware is used to log the request before passing it to the handler
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user email from context
		email, ok := r.Context().Value(authMiddleware.UserEmailKey).(string)
		if !ok || email == "" {
			email = "guest" // or handle it based on your requirements
		}

		// Capture IP address
		ip := getIPAddress(r)

		// Log the request to the database
		err := loggingService.LogToDB(email, r.URL.Path, ip)
		if err != nil {
			log.Printf("Failed to log user action: %v", err)
		}

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}
