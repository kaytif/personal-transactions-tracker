package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func transactionPostHandler(w http.ResponseWriter, r *http.Request) {
	// post is adding so no id check id will be generated automatically
	// name needs to be checked against the things
	// okay if transaction name is unique
	// you need the account id that links the transaction to the relevant account
	// need to add transactionn amount
	// category id as well

	// authenticate context received from r
	userID, ok := getUserID(r)
	if ok == false{
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
		}	

	var newTransaction Transaction

	// Here if there are any empty fields this error won't catch it. this error is for any json errors.
	err := json.NewDecoder(r.Body).Decode(&newTransaction)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// checked transaction name can not be empty
	if newTransaction.Name == "" {
		writeJSONError(w, "Name can not be empty", http.StatusBadRequest)
		return
	}

	// to deal with categories, we don't make them necessary
	// if category is not entered, then we want it to remain null
	// category remaining null means 

	// checked account affiliated must exist
	var verifyUserID int
	err = db.QueryRow(
		"SELECT id, user_id FROM accounts WHERE id = $1 AND deleted_at IS NULL",
		newTransaction.AccountID,
		//point is we check that hey does it exist and if it exists we just store it back
	).Scan(&newTransaction.AccountID, &verifyUserID)

	// give error if no active account with this ID exists.
	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, "Account not found, please create affiliated account first", http.StatusNotFound)
		return
	}

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
		}

	if verifyUserID != userID{
		writeJSONError(w, "Not Authorized", http.StatusForbidden)
		return
	}

	// we should check here if transaction amount is 0
	if newTransaction.Amount == 0 {
		writeJSONError(w, "Transaction amount can not be 0", http.StatusBadRequest)
		return
	}

	// check categoy affiliated must existed
	if newTransaction.CategoryID != nil {
		err = db.QueryRow(
		"SELECT id FROM categories WHERE id = $1",
		newTransaction.CategoryID,
		).Scan(&newTransaction.CategoryID)

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

	// Now let's make sure the date is all right
	// first we make sure that the transaction date is not its zero value

	// make sure it is not empty
	if newTransaction.Date == "" {
		writeJSONError(w, "Transaction date is required", http.StatusBadRequest)
    	return
	}

	// Parse the date
	transactionDate, err := time.Parse("2006-01-02", newTransaction.Date)
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

	// now inserting into postgre sql. we need to scan the rows and get the id
	// this is because the frontend might need to use it 
	err = db.QueryRow(
		"INSERT INTO transactions (account_id, name, amount, transaction_date, category_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		newTransaction.AccountID,
		newTransaction.Name,
		newTransaction.Amount,
		transactionDate,
		newTransaction.CategoryID,
	).Scan(&newTransaction.ID)

	if err != nil {
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Expense succesfull created
	w.WriteHeader(http.StatusCreated)

	// send created expense back as JSON
	json.NewEncoder(w).Encode(newTransaction)
}
