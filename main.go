package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Connect to PostgreSQL.
	connectDB()

	// Render provides PORT when deployed; use 8080 locally.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Transactions
	http.HandleFunc("/transactions", transactionPostHandler)

	// Users / authentication
	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/login", loginHandler)

	log.Println("server running on port " + port)

	// Start HTTP server.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("server failed: ", err)
	}
}