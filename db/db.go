package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var DB *sql.DB

func InitDB(dataSourceName string) *sql.DB {
	// Connection string format:
	// host=localhost port=5432 user=username password=password dbname=database_name sslmode=disable
	connStr := "host=localhost port=3080 user=user password=BVGbfHyDjxWAvkCaeYM4JU59ZnTt8p dbname=postgres sslmode=disable"

	// Establish a connection
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Unable to connect to the database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Unable to ping the database:", err)
	}

	log.Println("Connected to the PostgreSQL database successfully.")
	DB.SetMaxOpenConns(1)

	statement, err := DB.Prepare(`
    CREATE TABLE IF NOT EXISTS alerts (
        id SERIAL PRIMARY KEY,
        symbol VARCHAR(10) NOT NULL,
        trigger_value DECIMAL(10, 2) NOT NULL,
        alert_type TEXT CHECK (alert_type IN ('lower', 'higher')) NOT NULL,
        triggered BOOLEAN DEFAULT FALSE,
        user_id INT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        completed_at TIMESTAMP NULL DEFAULT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );
`)
	if err != nil {
		log.Fatal("Error preparing alerts table:", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error executing alerts table statement:", err)
	}

	statement, err = DB.Prepare(`
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
`)
	if err != nil {
		log.Fatal("Error preparing users table:", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error executing users table statement:", err)
	}

	statement, err = DB.Prepare(`
    CREATE TABLE IF NOT EXISTS logs (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255),
        endpoint VARCHAR(255),
        ip VARCHAR(45),
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
`)
	if err != nil {
		log.Fatal("Error preparing logs table:", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error executing logs table statement:", err)
	}

	if err != nil {
		log.Fatal("Error preparing logs table:", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal("Error executing logs table statement:", err)
	}

	DB.SetMaxOpenConns(50)
	DB.SetMaxIdleConns(50)
	DB.SetConnMaxLifetime(5 * time.Minute)

	return DB
}
