package main

import (
	"net/http"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"

)


func createCategory(w http.ResponseWriter, r *http.Request) {

	// get the verified user ID
	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// create categories variable
	var category Category

	// decode income request and store into varirable
	json.NewDecoder(r.Body).Decode(&category)

	// 
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// make sure name is not empty
	if category.Name == "" {
		writeJSONError(w, "Name can not be empty", http.StatusBadRequest)
		return
	}


	// insert in the relevant database
	err := db.QueryRow(
		`INSERT INTO categories (name, user_id) 
		VALUES ($1, $2)
		RETURNING id`,
		category.Name,
		userID,
	).Scan(&category.ID)


	// Duplicate error
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeJSONError(w, "Category name already exists", http.StatusConflict)
			return
		}

		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Category successfully created
	w.WriteHeader(http.StatusCreated)


	//Send the expense back as Json
	json.NewEncoder(w).Encode(category)

}