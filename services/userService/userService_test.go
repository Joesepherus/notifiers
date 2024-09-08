package userService

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	// Adjust this import path as needed
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	var err error
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	SetDB(db)

	// Run tests
	code := m.Run()

	// Clean up after tests
	// (if necessary, e.g., reset the database state)

	os.Exit(code)
}

func TestGetUserById_Success(t *testing.T) {
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
	mock.ExpectQuery("SELECT id, email FROM users")

	// Call the function you're testing
	users, err := GetUsers()

	// // Check if the function works as expected
	assert.Error(t, err)
	assert.Len(t, users, 0)
}

func TestGetUsers_ScanError(t *testing.T) {
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

func TestUpdatePassword_Success(t *testing.T) {
	// Mock the expected behavior
	mock.ExpectExec("UPDATE users SET password = ? WHERE email = ?").
		WithArgs(sqlmock.AnyArg(), "user@example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := UpdatePassword("user@example.com", "hashedPasswordValue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateUser_Success(t *testing.T) {
	// Define the query as used in the UpdatePassword function
	query := "INSERT INTO users (email, password) VALUES (?, ?)"

	// Mock Exec to return an error
	mock.ExpectExec(query).
		WithArgs("dushan@example.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the function you're testing
	userId, err := CreateUser("dushan@example.com", "hashedPasswordValue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// log.Print("err", err)
	// log.Print("userId", userId)
	assert.Equal(t, 1, userId)
}

func TestCreateUser_DBError(t *testing.T) {
	// Mock the expected query to return an error
	mock.ExpectExec("INSERT INTO users (email, password) VALUES (?, ?)").
		WithArgs("user@example.com", sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf("db exec error"))

	// Call the function you're testing
	userID, err := CreateUser("user@example.com", "password123")

	// Check the results
	assert.Equal(t, 0, userID)
	assert.EqualError(t, err, "failed to insert user: db exec error")
}
