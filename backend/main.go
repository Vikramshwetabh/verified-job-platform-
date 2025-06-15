// This is the entry point of our Go backend server for the Verified Job Platform
package main

import (
	"encoding/json"
	"fmt" // Used for printing messages to the console
	"log"
	"net/http" // Core package to build HTTP servers

	"database/sql"

	"github.com/gorilla/mux" // External package for advanced routing
	_ "github.com/lib/pq"
)

// creating a user struct
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var db *sql.DB // db is declare here to make its global and accessible in the route handler
const connStr = "host=localhost port=5432 user=postgres password=Vikram@32 dbname=verified_jobs sslmode=disable"

func main() {
	var err error                           // <-- declare only err here now
	db, err = sql.Open("postgres", connStr) // global db is assigned

	// Create a new router instance using Gorilla Mux
	router := mux.NewRouter()

	// Define a simple route for the home (root) URL "/"
	// This will respond with plain text when someone visits localhost:8080/
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// You should hash password here (we’ll add bcrypt later)

		_, err = db.Exec("INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)",
			user.Name, user.Email, user.Password, user.Role)

		if err != nil {
			http.Error(w, "Error saving user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "User registered successfully!")
	})

	// Log to the terminal that the server is live
	fmt.Println("Server running on http://localhost:8080")

	// my PostgreSQL details
	connStr := "host=localhost port=5432 user=postgres password=Vikram@32 dbname=verified_jobs sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(" Failed to connect to database:", err)
	}
	defer db.Close()

	fmt.Println("Connected to PostgreSQL database")

	// Start the HTTP server on port 8080 and use the router to handle requests
	http.ListenAndServe(":8080", router)
}
