package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Read ID from URL.
func deleteAccount(w http.ResponseWriter, r *http.Request) {

	idString := r.URL.Query().Get("id")

	// Convert from string to integer
	id, err := strconv.Atoi(idString)

	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Delete expense from PostgreSQL
	result, err := db.Exec(
		"UPDATE accounts SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL",
		id,
	)

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check whether if anything was deleted
	rowsAffected, err := result.RowsAffected()
	
	if err != nil {
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}


	// If not matching account existed
	if rowsAffected == 0 {
		writeJSONError(w, "Account not found", http.StatusBadRequest)
		return
	
	}

	// Successful delete no body return
	w.WriteHeader(http.StatusNoContent)
}


