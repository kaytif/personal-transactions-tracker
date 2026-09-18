package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

func createAccount(w http.ResponseWriter, r *http.Request) {

	// if err != nil || id <= 0 {
	// 	writeJSONError(w, "Invalid ID", http.StatusBadRequest)
	// 	return
	// }

	// Create an empty account
	var newAccount Account

	// decode JSON
	// if no id is entered i guess it is caught here

	err := json.NewDecoder(r.Body).Decode(&newAccount)

	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// now we have to check all the conditions needed for the new account to be right

	// first condition is that name can not be empty
	if newAccount.Name == "" {
		writeJSONError(w, "Name not given", http.StatusBadRequest)
		return
	}

	// Insert into db
	err = db.QueryRow(
		"INSERT INTO accounts (name) VALUES ($1) returning id",
		newAccount.Name,
	).Scan(&newAccount.ID) // we let postgre generate the new id and then we save it in our thing



	// Duplicate error
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeJSONError(w, "Account name already exists", http.StatusConflict)
			return
		}

		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}


		// Handle the opening balance
	if newAccount.Balance != 0{
		_, err = db.Exec(
			"INSERT INTO transactions (name, account_id, amount, transaction_date) VALUES ($1, $2, $3, $4)",
			newAccount.Name + ": Opening Balance",
			newAccount.ID,
			newAccount.Balance,
			time.Now(),
		)
	}

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Account successfully created
	w.WriteHeader(http.StatusCreated)

	// Send crated expense back as JSON
	json.NewEncoder(w).Encode(newAccount)
}

