package main

import (
	"encoding/json"
	"net/http"
)

func getAccount(w http.ResponseWriter, r *http.Request) {
	// Ask postgresql for all expenses
	rows, err := db.Query("SELECT id, name FROM accounts WHERE deleted_at IS NULL")
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Start with an empty slice to empty results return [] instead of null.
	accounts := []Account{}

	// Move through each row postgresq; returned
	for rows.Next() {
		var accountStore Account
		err := rows.Scan(&accountStore.ID, &accountStore.Name)
		if err != nil {
			writeJSONError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Add this account to our list.
		accounts = append(accounts, accountStore)
	}
	
	// Check whether any error occurred while iterating through any rows
	if err := rows.Err(); err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// send completed response back to client
	json.NewEncoder(w).Encode(accounts)
}

