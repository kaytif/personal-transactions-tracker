package main

import(
	"net/http"
	"encoding/json"
)

func createTransfer(w http.ResponseWriter, r *http.Request){

	// create a transfer var
	var transfer Transfer


	// decode request into transfer
	err := json.NewDecoder(r.Body).Decode(&transfer)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
	}

	// Now create the first transaction


	// Then create the second transaction

	// make two transactions



}