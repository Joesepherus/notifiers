package authUtils

import (
	"testing"
	"time"
	"tradingalerts/types/userTypes"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateToken(t *testing.T) {
	email := "user@example.com"
	token, err := GenerateToken(email)

	// Check if there's no error
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims := jwt.MapClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})
	assert.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok && parsedToken.Valid)

	// Verify the claims
	assert.Equal(t, email, claims["email"])
	expiration := time.Unix(int64(claims["exp"].(float64)), 0)
	assert.WithinDuration(t, time.Now().Add(time.Hour*24), expiration, time.Second*5)
}

func TestCheckPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	user := &userTypes.User{
		Password: string(hashedPassword),
	}

	// Test correct password
	assert.True(t, CheckPassword(user, password))

	// Test incorrect password
	assert.False(t, CheckPassword(user, "wrongpassword"))
}
