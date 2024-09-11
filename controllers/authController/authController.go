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
	"tradingalerts/services/loggingService"
	"tradingalerts/services/userService"
	"tradingalerts/utils/authUtils"
	"tradingalerts/utils/errorUtils"
	"tradingalerts/utils/subscriptionUtils"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	errorUtils.MethodNotAllowed_error(w, r)
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		log.Println("Error creating user: Email and password are required")
		loggingService.LogToDB("ERROR", "Error creating user: Email and password are required", r)
		http.Redirect(w, r, "/error?message=Email+and+password+are+required", http.StatusSeeOther)
		return
	}

	userID, err := userService.CreateUser(email, password)
	if err != nil {
		log.Println("Error creating user:", err)
		loggingService.LogToDB("ERROR", "Error creating user", r)
		http.Redirect(w, r, "/error?message=Error+creating+user", http.StatusSeeOther)
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
	errorUtils.MethodNotAllowed_error(w, r)

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		log.Println("Email and password are required")
		loggingService.LogToDB("ERROR", "Email and password are required", r)
		http.Redirect(w, r, "/error?message=Email+and+password+are+required", http.StatusSeeOther)
		return
	}

	user, err := userService.GetUserByEmail(email)
	if err != nil {
		log.Println("Invalid email or password")
		loggingService.LogToDB("ERROR", "Invalid email or password", r)
		http.Redirect(w, r, "/error?message=Invalid+email+or+password", http.StatusSeeOther)
		return
	}

	if !authUtils.CheckPassword(user, password) {
		log.Println("Invalid email or password")
		loggingService.LogToDB("ERROR", "Invalid email or password", r)
		http.Redirect(w, r, "/error?message=Invalid+email+or+password", http.StatusSeeOther)
		return
	}

	// Generate token
	token, err := authUtils.GenerateToken(email)
	if err != nil {
		log.Println("Error generating token")
		loggingService.LogToDB("ERROR", "Error generating token", r)
		http.Redirect(w, r, "/error?message=Error+generating+token", http.StatusSeeOther)
		return
	}

	// Set token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,                    // Prevent JavaScript access
		Secure:   true,                    // Use only on HTTPS
		SameSite: http.SameSiteStrictMode, // prevent CSRF
		MaxAge:   3600 * 24,               // Token expires in 1 day
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
	errorUtils.MethodNotAllowed_error(w, r)
	email := r.FormValue("email")

	if email == "" {
		log.Println("Email is required")
		loggingService.LogToDB("ERROR", "Email is required", r)
		http.Redirect(w, r, "/error?message=Email+is+required", http.StatusSeeOther)
		return
	}

	// Generate a random token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		log.Println("Error generating token")
		loggingService.LogToDB("ERROR", "Error generating token", r)
		http.Redirect(w, r, "/error?message=Error+generating+token", http.StatusSeeOther)
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
	errorUtils.MethodNotAllowed_error(w, r)
	token := r.FormValue("token")
	password := r.FormValue("password")

	tokenData, exists := authUtils.ResetTokens[token]
	if !exists {
		http.Redirect(w, r, "/token-expired", http.StatusSeeOther)
		return
	}

	// Check if the token has expired
	if time.Now().After(tokenData.Expiration) {
		log.Print("token has expired")
		delete(authUtils.ResetTokens, token)
		http.Redirect(w, r, "/token-expired", http.StatusSeeOther)
		return
	}
	log.Print("token is valid", tokenData.Expiration, time.Now())

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password")
		loggingService.LogToDB("ERROR", "Error hashing password", r)
		http.Redirect(w, r, "/error?message=Error+hashing+password", http.StatusSeeOther)
		return
	}

	// Save the new password in your database (pseudo-code)
	err = userService.UpdatePassword(tokenData.Email, string(hashedPassword))
	if err != nil {
		log.Println("Error saving new password")
		loggingService.LogToDB("ERROR", "Error saving new password", r)
		http.Redirect(w, r, "/error?message=Error+saving+new+password", http.StatusSeeOther)
		return
	}

	// Invalidate the token
	delete(authUtils.ResetTokens, token)
	http.Redirect(w, r, "/reset-password-success", http.StatusSeeOther)

	w.Write([]byte("Password has been reset successfully."))
}
