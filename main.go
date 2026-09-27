package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Connect to PostgreSQL.
	connectDB()

	// Use 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Transactions
	http.HandleFunc("/transactions", handlerTransaction)

	// Accounts
	http.HandleFunc("/accounts", handlerAccount)

	// Transfers
	http.HandleFunc("/transfers", handlerTransfer)

	// Categories
	http.HandleFunc("/categories", handlerCategory)

	// Users / authentication
	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/login", loginHandler)

	log.Println("server running on port " + port)

	// Start HTTP server.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("server failed: ", err)
	}
}
