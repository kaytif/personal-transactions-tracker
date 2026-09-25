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

	// Delete transaction from PostgreSQL
	result, err := db.Exec(
		`UPDATE transactions 								-- Modify the transaction table
		SET deleted_at = NOW() 								-- Soft deletion
		FROM accounts 										-- we want to make checks in accounts
		WHERE accounts.id = transactions.account_id 		-- match the right account with the transaction account
		AND accounts.user_id = $1 							-- match the right user with the user who made the account
		AND transactions.id = $2 							-- match the right transaction 
		AND transactions.deleted_at IS NULL 				-- transaction must not already be deleted`,
		userID,
		ID,
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





