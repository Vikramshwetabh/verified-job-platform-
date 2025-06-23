package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Job struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Company     string `json:"company"`
	PostedBy    string `json:"posted_by"` // recruiter email or id
}

// Pass DB instance from main.go
var DB *sql.DB

func SetDB(database *sql.DB) {
	DB = database
}

func GetJobs(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, title, description, company, posted_by FROM jobs")
	if err != nil {
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.Company, &job.PostedBy)
		if err != nil {
			http.Error(w, "Error reading jobs", http.StatusInternalServerError)
			return
		}
		jobs = append(jobs, job)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}
