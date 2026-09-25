package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

// The *http.request has a pointer because the handler receives the address
// instead of making a copy of the whole address
func transactionPutHandler(w http.ResponseWriter, r *http.Request){

	// we need to check that the context id we are receiving is integer and store it
	userID, ok := getUserID(r)
	if ok == false{
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// first thing to do is to make sure that the transaction id
	// you are receiving is not null
	idString := r.URL.Query().Get("id")

	// and does not have any issues
	// we do that by converting to string 
	// then we convert from string to integer
	id, err := strconv.Atoi(idString)
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var storeTransaction Transaction

	// if all is okay, then we decode the body into our struct
	err = json.NewDecoder(r.Body).Decode(&storeTransaction)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Let's check that the transaction actually exists based on the 
	// transaction id
	err = db.QueryRow(
		`SELECT id 
		FROM accounts
		JOIN transactions ON accounts.id = transactions.account_id
		WHERE accounts.user_id = $1
		AND transactions.deleted_at IS NULL
		`,
		userID,
	).Scan(&storeTransaction.ID)

	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, "Transaction does not exist", http.StatusNotFound)
		return
	}

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// once we do that, we need to do all the validations
	// the validations we need to do are, check that is is not something that is deleted
	
	// check transaction is 0
	// we should check here if transaction amount is 0
	if storeTransaction.Amount == 0 {
		writeJSONError(w, "Transaction amount can not be 0", http.StatusBadRequest)
		return
	}

	if storeTransaction.Name == "" {
		writeJSONError(w, "Transaction name is required", http.StatusBadRequest)
		return
	}

	// do all the date checks
	// make sure it is not empty
	if storeTransaction.Date == "" {
		writeJSONError(w, "Transaction date is required", http.StatusBadRequest)
    	return
	}

	// Parse the date
	transactionDate, err := time.Parse("2006-01-02", storeTransaction.Date)
	if err != nil {
		writeJSONError(w, "Invalid date format", http.StatusBadRequest)
    	return
	}

	// Make sure the date is not greater than now
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if transactionDate.After(today) {
		writeJSONError(w, "Transaction date can not be in the future", http.StatusBadRequest)
		return
	}

	// check categoy affiliated must existed
	if storeTransaction.CategoryID != nil {
		err = db.QueryRow(
		"SELECT id FROM categories WHERE id = $1",
		storeTransaction.CategoryID,
		).Scan(&storeTransaction.CategoryID)

		// give error if no active account with this ID exists.
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, "Category not found, please create affiliated category first", http.StatusNotFound)
			return
		}

		// give error if no active category with this ID exists.
		if err != nil {
			writeJSONError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

		// once all is said and done, let's update the transaction
	result, err := db.Exec(
		"UPDATE transactions SET name = $1, amount = $2, transaction_date = $3, updated_at = $4, category_id = $5 WHERE id = $6",
		storeTransaction.Name,
		storeTransaction.Amount,
		transactionDate,
		time.Now(),
		storeTransaction.CategoryID,
		id,
	)

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
		}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// no matching transaction existed
	if rowsAffected == 0 {
		writeJSONError(w, "Transaction nor found", http.StatusNotFound)
		return
	}

	// keep the ID from the url
	storeTransaction.ID = id

	// send updated transaction back as json
	json.NewEncoder(w).Encode(storeTransaction)
	}
	