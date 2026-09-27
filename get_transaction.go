package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func getTransaction(w http.ResponseWriter, r *http.Request) {

	// get the verified user ID
	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	pageString := r.URL.Query().Get("page")
	limitString := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageString)
	if err != nil || page <= 0 {
		writeJSONError(w, "Invalid page", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitString)
	if err != nil || limit <= 0 {
		writeJSONError(w, "Invalid limit", http.StatusBadRequest)
		return
	}

	offset := (page - 1) * limit

	// we need to make sure that we are getting the right transactions from the user id
	// so we neeed to create a transaction variable
	// create list that stores the transactions
	// get the transaction from postgres
	// iterate through and store in the list
	// decode and return
	// selects all the rowws
	rows, err := db.Query(
		`SELECT transactions.id, transactions.name, transactions.amount,
		transactions.transaction_date, transactions.created_at, accounts.name,
		categories.name, transactions.transaction_type
		FROM accounts
		JOIN transactions ON accounts.id = transactions.account_id
		LEFT JOIN categories ON transactions.category_id = categories.id
		WHERE accounts.user_id = $1
		AND transactions.deleted_at IS NULL
		ORDER BY transactions.transaction_date 
		DESC LIMIT $2 OFFSET $3`,
		userID,
		limit,
		offset,
	)

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Start with an empty slice so empty results return [] instead of null
	transactions := []TransactionResponse{}

	// Move through the rows
	for rows.Next() {
		var transactionResponse TransactionResponse

		err := rows.Scan(
			&transactionResponse.TransactionID,
			&transactionResponse.TransactionName,
			&transactionResponse.TransactionAmount,
			&transactionResponse.TransactionDate,
			&transactionResponse.TransactionCreatedAt,
			&transactionResponse.AccountName,
			&transactionResponse.CategoryName,
			&transactionResponse.TransactionType,
		)
		if err != nil {
			writeJSONError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// add transaction to our list

		transactions = append(transactions, transactionResponse)
	}

	// Check whether an error occurred while iterating through rows.
	if err := rows.Err(); err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// decoode and return
	json.NewEncoder(w).Encode(transactions)

}
