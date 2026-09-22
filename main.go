package main

import (
	"log"
	"net/http"
	"os"
)


// Remember the logic for main is to connect to database, then check connection, then start server
func main() {
	// Connect to PostgreSQL.
	connectDB()

	// Send API requests to handler

	// Render provides PORT when deployed, use 8080 locally
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// for accounts


	// register handler to routes for transactions
 	http.HandleFunc("/transactions", transactionPostHandler)

	
	// for categories
	


	// for users
	http.HandleFunc("/register", registerUser)


	log.Println("server running on port" + port)

	// Start the HTTP server,
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("server failed: ", err)
	}
}