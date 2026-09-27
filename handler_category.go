package main

import "net/http"


// expenseHandler sends each HTTP method to the correct handler
func handlerCategory(w http.ResponseWriter, r *http.Request) {
 	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
 	case http.MethodGet:
 		getCategory(w, r)

 	case http.MethodPost:
 		createCategory(w, r)

 	case http.MethodPut:
 		updateCategory(w, r)

 	case http.MethodDelete:
 		deleteCategory(w, r)
	
 	default:
 		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
 	}
}