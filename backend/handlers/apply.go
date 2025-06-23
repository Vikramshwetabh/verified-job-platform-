// handlers/apply.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Application struct {
	UserEmail string `json:"user_email"`
	Resume    string `json:"resume"` // Can be file path or base64
}

// ApplyToJob handles POST /jobs/{id}/apply
func ApplyToJob(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		jobID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid job ID", http.StatusBadRequest)
			return
		}

		var app Application
		if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Insert application into DB
		_, err = db.Exec(`
			INSERT INTO applications (user_email, job_id, resume)
			VALUES (?, ?, ?)`,
			app.UserEmail, jobID, app.Resume,
		)

		if err != nil {
			http.Error(w, "Failed to apply: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Application submitted successfully!")
	}
}
