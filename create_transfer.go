package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func createTransfer(w http.ResponseWriter, r *http.Request) {

	userID, ok := getUserID(r)
	if ok == false {
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// create a transfer var
	var transfer Transfer

	// decode request into transfer
	err := json.NewDecoder(r.Body).Decode(&transfer)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// we need to make sure that the transfer is not 0
	if transfer.Amount <= 0 {
		writeJSONError(w, "Amount can not be 0 or negative", http.StatusBadRequest)
		return
	}

	// make sure from and to account id is not the same
	if transfer.FromAccountID == transfer.ToAccountID {
		writeJSONError(w, "Can not transfer to same account", http.StatusBadRequest)
		return
	}

	// we need to do checks for the dates as well
	// make sure it is not empty
	if transfer.Date == "" {
		writeJSONError(w, "Transaction date is required", http.StatusBadRequest)
		return
	}

	// make sure it is not in the future

	// first parse the date
	var transferDate time.Time
	transferDate, err = time.Parse("2006-01-02", transfer.Date)
	if err != nil {
		writeJSONError(w, "Transaction date is required", http.StatusBadRequest)
		return
	}

	// Make sure the date is not greater than today
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if transferDate.After(today) {
		writeJSONError(w, "Transaction date can not be in the future", http.StatusBadRequest)
		return
	}

	// set outgoings and ingoings amount
	var outgoingAmount float64
	var incomingAmount float64
	outgoingAmount = -transfer.Amount
	incomingAmount = transfer.Amount

	// retrieve account name outgoing
	var outgoingAccountName string
	err = db.QueryRow(
		`SELECT name
		FROM accounts
		WHERE id =  $1
		AND user_id = $2
		AND deleted_at is NULL`,
		transfer.FromAccountID,
		userID,
	).Scan(&outgoingAccountName)

	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, "Account not found, please create affiliated account first", http.StatusNotFound)
		return
	}

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// retrieve account name incoming
	var incomingAccountName string
	err = db.QueryRow(
		`SELECT name
		FROM accounts
		WHERE id =  $1
		AND user_id = $2
		AND deleted_at is NULL`,
		transfer.ToAccountID,
		userID,
	).Scan(&incomingAccountName)

	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, "Account not found, please create affiliated account first", http.StatusNotFound)
		return
	}

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// give the transfer a name for later
	var transferName string
	transferName = "Transfer from " + outgoingAccountName + " to " + incomingAccountName

	// Now create the first transaction, say the outgoing transaction
	// remember that the process needs atomicity since either both
	// transaction gets created or none does

	tx, err := db.Begin()
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer tx.Rollback()

	var outgoingTransaction Transaction
	err = tx.QueryRow(
		`INSERT INTO transactions 
		(account_id, amount, name, transaction_type, transaction_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING account_id, amount, name, transaction_type, id`,
		transfer.FromAccountID,
		outgoingAmount,
		transferName,
		"transfer",
		transferDate,
	).Scan(
		&outgoingTransaction.AccountID,
		&outgoingTransaction.Amount,
		&outgoingTransaction.Name,
		&outgoingTransaction.TransactionType,
		&outgoingTransaction.ID,
	)

	// check for errors
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// make sure to the save the date
	outgoingTransaction.Date = transfer.Date

	// Then create the second transaction
	var incomingTransaction Transaction
	err = tx.QueryRow(
		`INSERT INTO 
		transactions (account_id, amount, name, transaction_type, transaction_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING account_id, amount, name, transaction_type, id`,
		transfer.ToAccountID,
		incomingAmount,
		transferName,
		"transfer",
		transferDate,
	).Scan(
		&incomingTransaction.AccountID,
		&incomingTransaction.Amount,
		&incomingTransaction.Name,
		&incomingTransaction.TransactionType,
		&incomingTransaction.ID,
	)

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// make sure to the save the date
	incomingTransaction.Date = transfer.Date

	err = tx.Commit()
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Account successfully created
	w.WriteHeader(http.StatusCreated)

	// Send transfer back as JSON
	json.NewEncoder(w).Encode(&transfer)

}
