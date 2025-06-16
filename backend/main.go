package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Global database connection variable
var db *sql.DB

// PostgreSQL connection string
const connStr = "host=localhost port=5432 user=postgres password=Vikram@32 dbname=verified_jobs sslmode=disable"

// User struct representing registration data
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Handler function to register a user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// NOTE: You should hash the password before saving (bcrypt recommended)

	_, err = db.Exec("INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)",
		user.Name, user.Email, user.Password, user.Role)

	if err != nil {
		http.Error(w, "Error saving user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully!")
}

func main() {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	fmt.Println("Connected to PostgreSQL database")

	// Create router and register routes
	router := mux.NewRouter()
	router.HandleFunc("/register", RegisterUser).Methods("POST")

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
