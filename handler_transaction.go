package main

import "net/http"

// expenseHandler sends each HTTP method to the correct handler
func handlerTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		getTransaction(w, r)

	case http.MethodPost:
		createTransaction(w, r)

	case http.MethodPut:
		updateTransaction(w, r)

	case http.MethodDelete:
		deleteTransaction(w, r)

	default:
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
