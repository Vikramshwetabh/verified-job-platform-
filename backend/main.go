// This is the entry point of our Go backend server for the Verified Job Platform
package main

import (
	"fmt" // Used for printing messages to the console
	"log"
	"net/http" // Core package to build HTTP servers

	"database/sql"

	"github.com/gorilla/mux" // External package for advanced routing
	_ "github.com/lib/pq"
)

func main() {
	// Create a new router instance using Gorilla Mux
	router := mux.NewRouter()

	// Define a simple route for the home (root) URL "/"
	// This will respond with plain text when someone visits localhost:8080/
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Write a response to the browser
		fmt.Fprintln(w, "Verified Job Platform API is running ")
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
