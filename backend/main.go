package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"./handlers"
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

	// Create jobs table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		company TEXT,
		location TEXT,
		recruiter_email TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("Failed to create jobs table:", err)
	}

	// Create applications table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_email TEXT NOT NULL,
		job_id INTEGER NOT NULL,
		status TEXT DEFAULT 'pending',
		resume TEXT,
		FOREIGN KEY (job_id) REFERENCES jobs(id)
	)`)
	if err != nil {
		log.Fatal("Failed to create applications table:", err)
	}

	fmt.Println("Connected to SQLite and tables are ready.")

	// Setup router and endpoints
	router := mux.NewRouter()

	// Pass DB to handlers
	handlers.SetDB(db)

	// Auth Routes
	router.HandleFunc("/register", handlers.RegisterHandler(db)).Methods("POST")
	router.HandleFunc("/login", handlers.LoginHandler(db)).Methods("POST")

	// Job Routes
	router.HandleFunc("/jobs", handlers.GetJobs).Methods("GET")
	router.HandleFunc("/jobs/{id}/apply", handlers.ApplyToJob(db)).Methods("POST")

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
