package userService

import (
	"errors"
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

func TestGetUsers_Success(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	SetDB(db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the expected query and the returned rows
	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow(1, "bob@gmail.com").
		AddRow(2, "dushan@gmail.com")

	mock.ExpectQuery("SELECT id, email FROM users").
		WillReturnRows(rows)

	// Call the function you're testing
	users, err := GetUsers()

	// // Check if the function works as expected
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "bob@gmail.com", users[0].Email)
	assert.Equal(t, "dushan@gmail.com", users[1].Email)
}

func TestGetUsers_NoUsers(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	SetDB(db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, email FROM users")

	// Call the function you're testing
	users, err := GetUsers()

	// // Check if the function works as expected
	assert.Error(t, err)
	assert.Len(t, users, 0)
}

func TestGetUsers_ScanError(t *testing.T) {
	// Create mock DB and mock query results
	db, mock, err := sqlmock.New()
	SetDB(db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the expected query and a faulty row (mismatching columns or types)
	rows := sqlmock.NewRows([]string{"id", "email"}).
		AddRow("invalid", "dushan@gmail.com")

	mock.ExpectQuery("SELECT id, email FROM users").
		WillReturnRows(rows)

	// Call the function you're testing
	users, err := GetUsers()

	// Check if the scanning error is handled correctly
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), "failed to scan row")
}

func TestGetUsers_QueryError(t *testing.T) {
	// Create mock DB and mock query results
	db, mock, err := sqlmock.New()
	SetDB(db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the expected query to return an error
	mock.ExpectQuery("SELECT id, email FROM users").
		WillReturnError(errors.New("query error"))

	// Call the function you're testing
	users, err := GetUsers()

	// Check if the error is handled correctly
	assert.Nil(t, users)
	assert.EqualError(t, err, "failed to query users: query error")
}
