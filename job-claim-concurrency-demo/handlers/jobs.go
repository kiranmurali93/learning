package handlers

import (
	"encoding/json"
	"job-claim-concurrency-demo/db"
	"net/http"
)

func ClaimJob(w http.ResponseWriter, r *http.Request) {
	jobId := r.URL.Query().Get("jobId")
	recruiterId := r.URL.Query().Get("recruiterId")

	// usually recruiter id will be taken from the auth token
	if jobId == "" || recruiterId == "" {
		http.Error(w, "Missing jobId or recruiterId", http.StatusBadRequest)
		return
	}

	// check if the job is eligible for claiming
	tx, err := db.DB.Begin()

	if err != nil {
		http.Error(w, "Could not start the transaction", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback()

	// check if already not claimed
	var exists bool
	check := `SELECT EXISTS(SELECT 1 FROM job_claims where job_id = $1 and recruiter_id = $2)`
	err = tx.QueryRow(check, jobId, recruiterId).Scan(&exists)
	if err != nil {
		http.Error(w, "Could not check the status", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Recuriter already claimed this job", http.StatusConflict)
		return
	}

	// check for count should not be greater than 4
	var count int
	query := `SELECT count FROM job_claims_count WHERE job_id = $1 FOR UPDATE;`
	err = tx.QueryRow(query, jobId).Scan(&count)

	if err != nil {
		http.Error(w, "Failed to check the count", http.StatusInternalServerError)
		return
	}

	// here 4 is the maximum number of job claimes that can happen
	if count >= 4 {
		http.Error(w, "Job claimed fully", http.StatusForbidden)
		return
	}

	// insert to the table
	insert := `INSERT INTO job_claims (job_id, recruiter_id) VALUES ($1, $2)`
	_, err = tx.Exec(insert, jobId, recruiterId)
	if err != nil {
		http.Error(w, "Failed to claime the job", http.StatusInternalServerError)
		return
	}

	// update the counter
	updateCounter := `INSERT INTO job_claims_count (job_id, count) VALUES ($1, 1) ON CONFLICT (job_id) DO UPDATE SET count = job_claims_count.count + 1;`
	_, err = tx.Exec(updateCounter, jobId)

	if err != nil {
		http.Error(w, "Failed to update claims count", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()

	if err != nil {
		http.Error(w, "Failed to commit claim", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Job Claimed.",
	})
}
