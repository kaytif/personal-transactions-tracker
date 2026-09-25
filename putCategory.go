package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"github.com/jackc/pgx/v5/pgconn"
)

func updateCategory(w http.ResponseWriter, r *http.Request){
	
	
	// first receive string
	idString := r.URL.Query().Get("id")	

	// then convert string into integer
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

	var category Category

	// decode what you get and store it
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// validation checks on name, make sure that the name is not empty
	if category.Name == ""{
		writeJSONError(w, "Name can not be empty", http.StatusBadRequest)
		return
	}

	// update whatever
	result, err := db.Exec(
		`UPDATE categories SET name = $1 WHERE id= $2 AND user_id = $3 AND deleted_at IS NULL`,
		category.Name,
		id,
		userID,
		)

	if err != nil{
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeJSONError(w, "Category name already exists", http.StatusConflict)
			return
		}
		
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}


	rowsAffected, err := result.RowsAffected()
	
	if err != nil {
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return		
	}

	if rowsAffected == 0 {
		writeJSONError(w, "Category not found", http.StatusNotFound)
		return		
	}

	// keep the ID from the URL 
	category.ID = id

	// send the updated category back as json
	json.NewEncoder(w).Encode(category)
}