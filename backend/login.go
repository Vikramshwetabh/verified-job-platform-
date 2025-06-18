// login.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// LoginRequest struct represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler handles user login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	var role string
	query := "SELECT password, role FROM users WHERE email = ?"
	err := db.QueryRow(query, loginReq.Email).Scan(&hashedPassword, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Compare the hashed password with the provided one
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginReq.Password)) != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Send success response
	response := fmt.Sprintf("Login successful! Role: %s", role)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, response)
}
