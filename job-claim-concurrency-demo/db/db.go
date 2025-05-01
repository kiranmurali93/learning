package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")

	var err error
	DB, err = sql.Open("postgres", DB_CONNECTION_STRING)
	if err != nil {
		log.Fatal("Error connecting DB", err)
	}


    if err = DB.Ping(); err != nil {
		log.Fatal("Could not ping the database", err)
    }

	log.Printf("Database connected successfully.")
}