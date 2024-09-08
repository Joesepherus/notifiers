package authController

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"tradingalerts/mail"
	"tradingalerts/payments/payments"
	"tradingalerts/services/userService"
	"tradingalerts/utils/authUtils"
	"tradingalerts/utils/subscriptionUtils"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	userID, err := userService.CreateUser(email, password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	payments.CreateCustomer(email)
	canAddAlert, subscriptionType := subscriptionUtils.CheckToAddAlert(userID, email)
	subscriptionUtils.UserSubscription[email] = subscriptionUtils.UserAlertInfo{
		CanAddAlert:      canAddAlert,
		SubscriptionType: subscriptionType,
	}
	http.Redirect(w, r, "/?login=true", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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

	if !authUtils.CheckPassword(user, password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := authUtils.GenerateToken(email)
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

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	email := r.FormValue("email")

	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Generate a random token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Set the expiration time (e.g., 24 hours from now)
	expiration := time.Now().Add(24 * time.Hour)

	// Store the token with email and expiration time
	authUtils.ResetTokens[token] = authUtils.ResetTokenData{
		Email:      email,
		Expiration: expiration,
	}

	// Send the reset link via email
	resetLink := fmt.Sprintf(os.Getenv("URL")+"?token=%s", token)
	go mail.SendEmail(email, "Trading Alerts: Password Reset", fmt.Sprintf(
		"Click the link below to reset your password:"+resetLink,
	))

	http.Redirect(w, r, "/reset-password-sent", http.StatusSeeOther)

	w.Write([]byte("Password reset email sent."))
}

func SetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	token := r.FormValue("token")
	password := r.FormValue("password")

	tokenData, exists := authUtils.ResetTokens[token]
	if !exists {
		http.Redirect(w, r, "/token-expired", http.StatusSeeOther)
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	// Check if the token has expired
	if time.Now().After(tokenData.Expiration) {
		log.Print("token has expired")
		delete(authUtils.ResetTokens, token)
		http.Redirect(w, r, "/token-expired", http.StatusSeeOther)
		http.Error(w, "Token has expired", http.StatusBadRequest)
		return
	}
	log.Print("token is valid", tokenData.Expiration, time.Now())

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Save the new password in your database (pseudo-code)
	err = userService.UpdatePassword(tokenData.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error saving new password", http.StatusInternalServerError)
		return
	}

	// Invalidate the token
	delete(authUtils.ResetTokens, token)
	http.Redirect(w, r, "/reset-password-sucess", http.StatusSeeOther)

	w.Write([]byte("Password has been reset successfully."))
}
