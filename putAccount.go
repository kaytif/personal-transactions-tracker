package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"github.com/jackc/pgx/v5/pgconn"
	"database/sql"
)

// a json decoding error for example occurs when something that needs to be integer in struct is entered as something else in the input
// Put updates, and doing it first just for accounts
func updateAccount(w http.ResponseWriter, r *http.Request) {
	// Read ID from URL
	idString := r.URL.Query().Get("id")

	// Convert ID from string into integer
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

	// Create an empty account
	var newAccount Account

	// Decode incoming JSON.
	err = json.NewDecoder(r.Body).Decode(&newAccount)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate account name is not empty
	if newAccount.Name == "" {
		writeJSONError(w, "Account name is required", http.StatusBadRequest)
		return
	}

	// in soft delete: row remains in DB and deleted at timestamp marks it as deleted
	// If new account does not exist
	
	// first we need to find whether the account already exists, so we will try to retrive using SELECT

	// note that we don't use rows affected here because rows affected occurs as a result of changing data
	var checkAccount Account
	
	err = db.QueryRow(
		"SELECT id, name, deleted_at FROM accounts WHERE id = $1 and user_id = $2",
		id,
		userID,
	).Scan(&checkAccount.ID, &checkAccount.Name, &checkAccount.DeletedAt)


	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, "Account not found", http.StatusNotFound)
		return
	}

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if checkAccount.DeletedAt != nil {
		writeJSONError(w, "Account was deleted", http.StatusBadRequest)
		return
	}


	// Validate account name is unique
		// postgre sql does it for you, here we show the error if postgresql gives us an error
		// Note that you can only update the new account name and not id because id is already assigned itself by sql
	result, err := db.Exec(
		"UPDATE accounts SET name = $1 WHERE id = $2 and user_id = $3", 
		newAccount.Name,
		id,
		userID,
	)


	if err != nil {
		var pgErr *pgconn.PgError
		
		// Duplicate account name
		// PostgreSQL code 23505 = UNIQUE constraint violation
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeJSONError(w, "Account name already exists", http.StatusConflict)
			return
		}
	// We check here for an error in case there was an error in updating expenses
		// Here the error checks for errors that happen while PostgreSQL is executing the UPDATE for e.g.,;
		// database connection died, unique contraint violated, check ocnstraint violated, foreing key viollated, sql statement itself is invalid

		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Check how many rows were updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// No matching expense existed
	if rowsAffected == 0 {
		writeJSONError(w, "Account not found", http.StatusNotFound)
		return
	}

	// Keep the ID from the URL
	newAccount.ID = id

	// Send the updated account back as JSON
	json.NewEncoder(w).Encode(newAccount)

}