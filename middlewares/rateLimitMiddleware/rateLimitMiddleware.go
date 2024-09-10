package rateLimitMiddleware

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
	clients = make(map[string]*clientData)
	mu      sync.Mutex
)

type clientData struct {
    limiter *rate.Limiter
    banUntil time.Time
}

func getClientLimiter(ip string) *clientData {
	mu.Lock()
	defer mu.Unlock()

    log.Println("clients", clients)
	if client, exists := clients[ip]; exists {
		return client
	}

	limiter := rate.NewLimiter(3, 5) // 3 request per second, burst of 5
    client := &clientData {
        limiter: limiter,
        banUntil: time.Time{},
    }
	clients[ip] = client

	// Optionally, clean up old limiters after some time
	go func() {
		time.Sleep(10 * time.Minute)
		mu.Lock()
		delete(clients, ip)
		mu.Unlock()
	}()

	return client
}

func RateLimitPerClient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email, ok := r.Context().Value(authMiddleware.UserEmailKey).(string)
		if !ok || email == "" {
			email = "guest" // or handle it based on your requirements
		}
		ip := authUtils.GetIPAddress(r)

		client := getClientLimiter(ip)

        if time.Now().Before(client.banUntil) {
            http.Error(w, "You are temporarily banned. Please try again later.", http.StatusMethodNotAllowed)
            return
        }

		if !client.limiter.Allow() {
			err := loggingService.LogToDB(email, r.URL.Path+" - Too many requests from your IP", ip)
			if err != nil {
				log.Printf("Failed to log user action: %v", err)
			}

            mu.Lock()
            client.banUntil = time.Now().Add(5 * time.Minute)
            mu.Unlock()

			http.Error(w, "Too many requests from your IP, please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
