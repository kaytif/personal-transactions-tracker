package main

import (
	"net/http"
	"encoding/json"
)

func getTransfer(w http.ResponseWriter, r *http.Request){

	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Select the rows that are needed
	rows, err := db.Query(
		`SELECT transactions.id, 
				transactions.name, 
				transactions.amount, 
				transactions.transaction_date, 
				transactions.created_at, 
				accounts.name, 
				categories.name, 
				transactions.transaction_type
		FROM accounts
		JOIN transactions ON accounts.id = transactions.account_id
		LEFT JOIN categories ON transactions.category_id = categories.id
		WHERE transactions.transaction_type = $1
		AND accounts.user_id = $2
		AND transactions.deleted_at IS NULL`,
		"transfer",
		userID,
	)
		if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

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

		transactions = append(transactions, transactionResponse )
	}
	
	// Check whether an error occurred while iterating through rows.
	if err := rows.Err(); err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// decoode and return
	json.NewEncoder(w).Encode(transactions)

}
