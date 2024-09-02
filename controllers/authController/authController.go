package authController

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"notifiers/mail"
	"notifiers/services/userService"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}
		log.Printf("email, pass", email, password)
		_, err := userService.CreateUser(email, password)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// GenerateToken generates a JWT token
func GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		user, err := userService.GetUserByEmail(email)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		if err != nil || !userService.CheckPassword(user, password) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Generate token
		token, err := GenerateToken(email)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		// Set token in a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,      // Prevent JavaScript access
			Secure:   true,      // Use only on HTTPS
			MaxAge:   3600 * 24, // Token expires in 1 day
		})

		// Redirect to user dashboard or home page after successful login
		http.Redirect(w, r, "/alerts", http.StatusSeeOther)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the authentication token by setting an expired cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Setting MaxAge to -1 deletes the cookie
	})

	// Optionally, you can also invalidate the session or token on the server-side

	// Redirect to the homepage or login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

var resetTokens = map[string]string{}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	// Generate a random token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Store the token with an expiration time
	resetTokens[token] = email

	// Send the reset link via email
	resetLink := fmt.Sprintf(os.Getenv("URL")+"?token=%s", token)
	go mail.SendEmail(email, "Trading Alerts: Password Reset", fmt.Sprintf(
		"Click the link below to reset your password:"+resetLink,
	))
	if err != nil {
		http.Error(w, "Error sending email", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/reset-password-sent", http.StatusSeeOther)

	w.Write([]byte("Password reset email sent."))
}

func SetPassword(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	password := r.FormValue("password")
	log.Printf("password", password)

	email, ok := resetTokens[token]
	if !ok {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Save the new password in your database (pseudo-code)
	err = userService.UpdatePassword(email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error saving new password", http.StatusInternalServerError)
		return
	}

	// Invalidate the token
	delete(resetTokens, token)
	http.Redirect(w, r, "/reset-password-sucess", http.StatusSeeOther)

	w.Write([]byte("Password has been reset successfully."))
}
