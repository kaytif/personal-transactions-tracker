package main

import "net/http"

// expenseHandler sends each HTTP method to the correct handler
func handlerTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		getTransfer(w, r)

	case http.MethodPost:
		createTransfer(w, r)

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
