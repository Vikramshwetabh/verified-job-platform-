package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Vikramshwetabh/verified-job-platform/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./verified_jobs.db")
	if err != nil {
		log.Fatal("Failed to connect to SQLite:", err)
	}
	defer db.Close()

	// Create users table if it doesn't exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	fmt.Println("Connected to SQLite and users table ready.")

	router := mux.NewRouter()
	router.HandleFunc("/register", handlers.RegisterHandler(db)).Methods("POST")
	router.HandleFunc("/login", handlers.LoginHandler(db)).Methods("POST")

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// This code initializes a simple HTTP server with user registration and login functionality using SQLite.
// It creates a users table if it doesn't exist and sets up routes for registration and login.
