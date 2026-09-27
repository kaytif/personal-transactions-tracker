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

	http.Handle("/transactions", authMiddleware(http.HandlerFunc(handlerTransaction)))
	http.Handle("/accounts", authMiddleware(http.HandlerFunc(handlerAccount)))
	http.Handle("/transfers", authMiddleware(http.HandlerFunc(handlerTransfer)))
	http.Handle("/categories", authMiddleware(http.HandlerFunc(handlerCategory)))

	// Users / authentication
	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/login", loginHandler)

	log.Println("server running on port " + port)

	// Start HTTP server.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("server failed: ", err)
	}
}
