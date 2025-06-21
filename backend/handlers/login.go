package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// LoginRequest holds login credentials
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User struct holds user registration info
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// RegisterHandler handles user registration
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)",
			user.Name, user.Email, user.Password, user.Role)
		if err != nil {
			http.Error(w, "Failed to register: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "User registered successfully!")
	}
}

// LoginHandler handles user login
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var dbPassword, role string
		err := db.QueryRow("SELECT password, role FROM users WHERE email = ?", req.Email).Scan(&dbPassword, &role)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		if req.Password != dbPassword {
			http.Error(w, "Incorrect password", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Login successful! Role: %s", role)
	}
}
