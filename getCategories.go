package main

import (
	"encoding/json"
	"net/http"
)



func getCategories(w http.ResponseWriter, r *http.Request){
	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// first we want the sql query to get all the categories
	rows, err := db.Query(
		`SELECT name, id FROM categories WHERE user_id = $1 AND deleted_at IS NULL`,
		userID,
	)

	if err != nil {
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	categories := []Category{}
	
	for rows.Next(){

		var category Category
		err = rows.Scan(
			&category.Name,
			&category.ID,
		)

		if err != nil {
			writeJSONError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)

	}

		// Check whether any error occurred while iterating through any rows
	if err := rows.Err(); err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// send completed response back to client
	json.NewEncoder(w).Encode(categories)
}
