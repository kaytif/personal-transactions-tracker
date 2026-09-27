package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func getAccountByID(w http.ResponseWriter, r *http.Request) {
	// Read ID from URL.
	idString := r.URL.Query().Get("id")

	// Convert ID from string to integer.
	id, err := strconv.Atoi(idString)
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// get the verified user ID
	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Store the account returned by PostgreSQL.
	var accountStore Account

	// Retrieve one active account.
	err = db.QueryRow(
		`SELECT id, name 
		FROM accounts 
		WHERE id = $1 
		AND user_id = $2
		AND deleted_at IS NULL`,
		id, 
		userID,
	).Scan(&accountStore.ID, &accountStore.Name)

	// No active account with this ID exists.
	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, "Account not found", http.StatusNotFound)
		return
	}

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	balance, err := getBalance(accountStore.ID)
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	accountStore.Balance = balance

	// Send account back as JSON.
	json.NewEncoder(w).Encode(accountStore)
}