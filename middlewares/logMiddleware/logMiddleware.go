package logMiddleware

import (
	"net/http"
	"tradingalerts/middlewares/authMiddleware"
	"tradingalerts/services/loggingService"
	"tradingalerts/utils/authUtils"
)

// LogMiddleware is used to log the request before passing it to the handler
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user email from context
		email, ok := r.Context().Value(authMiddleware.UserEmailKey).(string)
		if !ok || email == "" {
			email = "guest" // or handle it based on your requirements
		}

		// Capture IP address
		ip := authUtils.GetIPAddress(r)

		// Log the request to the database
		loggingService.LogToDB(email, "INFO", "Accessing page", r.URL.Path, ip)

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}
