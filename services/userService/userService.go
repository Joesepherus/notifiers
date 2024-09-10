package userService

import (
	"database/sql"
	"fmt"
	"log"
	"tradingalerts/types/userTypes"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func CreateUser(email, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed hashing password: %v", err)
	}
	result, err := db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, string(hashedPassword))
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	log.Printf("User created with ID: %d and email: %s", userID, email)
	return int(userID), nil
}

func GetUserById(id int) (*userTypes.User, error) {
	user := &userTypes.User{}
	err := db.QueryRow("SELECT id, email, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %v", err)
	}
	return user, nil
}

func GetUserByEmail(email string) (*userTypes.User, error) {
	user := &userTypes.User{}
	err := db.QueryRow("SELECT id, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %v", err)
	}
	return user, nil
}

func GetUsers() ([]*userTypes.User, error) {
	// Prepare the query to select all users
	rows, err := db.Query("SELECT id, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	defer rows.Close()

	var users []*userTypes.User

	// Iterate through the rows
	for rows.Next() {
		user := &userTypes.User{}
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func UpdatePassword(email string, hashedPassword string) error {
	query := `UPDATE users SET password = $1 WHERE email = $2`

	_, err := db.Exec(query, hashedPassword, email)
	if err != nil {
		return fmt.Errorf("failed to update row: %v", err)
	}

	return nil
}
