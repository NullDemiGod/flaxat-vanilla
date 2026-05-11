package db

import (
	"log"
	"os"
	"database/sql"

	// Postgresql driver for golang, the underscore means import for side effects only
	// Meaning, we don't call anything directly
	_"github.com/lib/pq"
)


var DB *sql.DB

func Connect() {
	// First, get the connection string from .env file
	connectionString := os.Getenv("DATABASE_URL")

	// Second, we open the database using the connection string we just got and the postgres indirect import
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Failed to open database", err)
	}

	// Use Ping() to verify connection availability
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	log.Println("Connected to database successfuly")
}
