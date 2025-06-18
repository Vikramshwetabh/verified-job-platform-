package main

import (
	"database/sql"  // Standard library to work with SQL databases
	"encoding/json" // For decoding JSON requests
	"fmt"           // For formatted I/O
	"log"           // For logging errors/info
	"net/http"      // For HTTP server

	"github.com/gorilla/mux"        // Router package for better route handling
	_ "github.com/mattn/go-sqlite3" // SQLite driver (used as a blank import to register the driver)
)

// Declare a global database connection variable
var db *sql.DB

// Define a struct to hold user registration data
type User struct {
	Name     string `json:"name"`     // Name of the user
	Email    string `json:"email"`    // Email of the user
	Password string `json:"password"` // Password (should be hashed in production)
	Role     string `json:"role"`     // Role like job_seeker, recruiter etc.
}

// Handler function to register a new user via POST request
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body and decode JSON into User struct
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Insert user data into the users table
	_, err = db.Exec(`INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)`,
		user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		http.Error(w, "Error saving user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully!")
}

// Entry point of the application
func main() {
	var err error

	// Open a connection to SQLite database file (creates file if not exists)
	db, err = sql.Open("sqlite3", "./verified_jobs.db")
	if err != nil {
		log.Fatal("Failed to connect to SQLite database:", err)
	}
	defer db.Close()

	// Create users table if it does not exist
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

	// Create a router using gorilla/mux
	router := mux.NewRouter()

	// Register route for /register using POST method
	router.HandleFunc("/register", RegisterUser).Methods("POST") // RegisterUser function defined above
	router.HandleFunc("/login", LoginHandler).Methods("POST")    // Assuming LoginHandler is defined in login.go

	// Start the HTTP server on port 8080
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
