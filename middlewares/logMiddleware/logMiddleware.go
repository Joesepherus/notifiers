package logMiddleware

import (
	"net/http"
	"tradingalerts/services/loggingService"
)

// LogMiddleware is used to log the request before passing it to the handler
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request to the database
		loggingService.LogToDB("INFO", "Accessing page", r)

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}
