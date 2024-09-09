package authUtils

import (
	"net/http"
	"strings"
	"time"
	"tradingalerts/types/userTypes"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type ResetTokenData struct {
	Email      string
	Expiration time.Time
}

var ResetTokens = map[string]ResetTokenData{}

// GenerateToken generates a JWT token
func GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key"))
}

func CheckPassword(user *userTypes.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// Helper function to get the client's real IP address, including proxies
func GetIPAddress(r *http.Request) string {
	ip := r.RemoteAddr
	// Check if the request is coming from a proxy and get the real client IP
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		ip = strings.Split(forwardedFor, ",")[0] // Get the first IP in the chain
	}
	return ip
}
