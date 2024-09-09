package ratelimitmiddleware

import (
	"log"
	"net/http"
	"sync"
	"time"
	"tradingalerts/middlewares/authMiddleware"
	"tradingalerts/services/loggingService"
	"tradingalerts/utils/authUtils"

	"golang.org/x/time/rate"
)

var (
	clients = make(map[string]*rate.Limiter)
	mu      sync.Mutex
)

func getClientLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if limiter, exists := clients[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(1, 3) // 1 request per second, burst of 3
	clients[ip] = limiter

	// Optionally, clean up old limiters after some time
	go func() {
		time.Sleep(10 * time.Minute)
		mu.Lock()
		delete(clients, ip)
		mu.Unlock()
	}()

	return limiter
}

func RateLimitPerClient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, ok := r.Context().Value(authMiddleware.UserEmailKey).(string)
		if !ok || email == "" {
			email = "guest" // or handle it based on your requirements
		}
		ip := authUtils.GetIPAddress(r)

		limiter := getClientLimiter(ip)
		if !limiter.Allow() {
			err := loggingService.LogToDB(email, r.URL.Path, ip)
			if err != nil {
				log.Printf("Failed to log user action: %v", err)
			}
			http.Error(w, "Too many requests from your IP, please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
