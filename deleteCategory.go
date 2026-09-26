package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func deleteCategory(w http.ResponseWriter, r *http.Request){
	// extract id from url and into string
	idString := r.URL.Query().Get("id")

	// Convert from string to integer
	id, err := strconv.Atoi(idString)

	// get the verified user ID
	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		`UPDATE categories
		SET deleted_at = NOW()
		WHERE id = $1
		AND user_id = $2
		AND deleted_at IS NULL`, 
		id,
		userID,
	)

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// if matching account does not exist
	if rowsAffected == 0 {
		writeJSONError(w, "Category not found", http.StatusNotFound)
		return
	}

	// successful delete no body return
	w.WriteHeader(http.StatusNoContent)

}