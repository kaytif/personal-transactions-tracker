package main

import (
	"net/http"
	"strconv"
)

// Read ID from URL.
func deleteTransaction(w http.ResponseWriter, r *http.Request) {

	idString := r.URL.Query().Get("id")

	// Convert from string to integer
	id, err := strconv.Atoi(idString)

	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Delete transaction from PostgreSQL
	result, err := db.Exec(
		"UPDATE transactions SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL",
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


	// If not matching account not existed
	if rowsAffected == 0 {
		writeJSONError(w, "Transaction not found", http.StatusNotFound)
		return
	
	}

	// Successful delete no body return
	w.WriteHeader(http.StatusNoContent)
}





