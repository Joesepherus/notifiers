package userService

import (
	"errors"
	"log"
	"testing"

	// Adjust this import path as needed
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGetUserById_Success(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	SetDB(db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("hello123"), bcrypt.DefaultCost)

	// Mock the expected query and the returned rows
	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, "bob@gmail.com", hashedPassword)

	mock.ExpectQuery("SELECT id, email, password FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	// Call the function you're testing
	userById1, err := GetUserById(1)
	log.Print("userById1", userById1)
	if err != nil {
		t.Fatalf("unexpected error when calling GetUserById: %v", err)
	}
	// // Check if the function works as expected
	assert.NoError(t, err)
	assert.Equal(t, 1, userById1.ID)
	assert.Equal(t, "bob@gmail.com", userById1.Email)
	assert.Equal(t, string(hashedPassword), userById1.Password)
}

func TestGetUserById_Fail(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	SetDB(db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, email, password FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnError(errors.New("query error"))

	// Call the function you're testing
	userById1, err := GetUserById(1)

	// // Check if the function works as expected
	assert.Nil(t, userById1)
	assert.EqualError(t, err, "failed to query user: query error")
}

func TestGetUserByEmail_Success(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	SetDB(db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("hello123"), bcrypt.DefaultCost)

	// Mock the expected query and the returned rows
	rows := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, "bob@gmail.com", hashedPassword)

	mock.ExpectQuery("SELECT id, email, password FROM users WHERE email = ?").
		WithArgs("bob@gmail.com").
		WillReturnRows(rows)

	// Call the function you're testing
	userByEmail, err := GetUserByEmail("bob@gmail.com")
	log.Print("userByEmail", userByEmail)

	// // Check if the function works as expected
	assert.NoError(t, err)
	assert.Equal(t, 1, userByEmail.ID)
	assert.Equal(t, "bob@gmail.com", userByEmail.Email)
	assert.Equal(t, string(hashedPassword), userByEmail.Password)
}

func TestGetUserByEmail_Fail(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	SetDB(db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, email, password FROM users WHERE email = ?").
		WithArgs("bob@gmail.com").
		WillReturnError(errors.New("query error"))

	// Call the function you're testing
	userByEmail, err := GetUserByEmail("bob@gmail.com")

	// // Check if the function works as expected
	assert.Nil(t, userByEmail)
	assert.EqualError(t, err, "failed to query user: query error")
}
