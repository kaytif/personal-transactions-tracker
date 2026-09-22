package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"database/sql"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {

	// so essentially you need to check tht the username exists in the database
	// we need to retrive the password hash

	var newLogin Login
	var storedUser User
	// start by decoding the json
	err := json.NewDecoder(r.Body).Decode(&newLogin)
	if err != nil {
		writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}


	// we gotta check if the username exists
	err = db.QueryRow(
		"SELECT username, password_hash FROM users WHERE username = $1",
		newLogin.Username,	
	).Scan(&storedUser.Username, &storedUser.PasswordHash)
	
	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// now we gotta check if the password is right
	// for this we use some bcrypt var and functions
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(newLogin.Password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		writeJSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err != nil{
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
