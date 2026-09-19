package main

import "net/http"


// expenseHandler sends each HTTP method to the correct handler
func accountHandler(w http.ResponseWriter, r *http.Request) {
 	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
 	case http.MethodGet:
 		getAccount(w, r)

 	case http.MethodPost:
 		createAccount(w, r)

 	case http.MethodPut:
 		updateAccount(w, r)

 	case http.MethodDelete:
 		deleteAccount(w, r)
	
 	default:
 		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
 	}
}