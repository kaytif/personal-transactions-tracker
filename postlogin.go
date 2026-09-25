package main

import (
	"context"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string
const userIDKey contextKey = "userID"

// Get the authenticated user's ID from the request context.
func getUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(userIDKey).(int)
	return userID, ok
}

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
		"SELECT id, username, password_hash FROM users WHERE username = $1",
		newLogin.Username,
	).Scan(&storedUser.ID, &storedUser.Username, &storedUser.PasswordHash)

	if errors.Is(err, sql.ErrNoRows) {
		writeJSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Printf("failed to query user", err)
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

	if err != nil {
		log.Printf("failed to compre password hard: %v", err)
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// create a session
	var newSession Session
	var token string
		// create token using crypto/rand
	token = rand.Text()

	// Insert the session into postgresql
	err = db.QueryRow(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING user_id, token, created_at, expires_at",
		storedUser.ID,
		token,
		time.Now().Add(7 * 24 * time.Hour),
	).Scan(&newSession.UserID, &newSession.Token, &newSession.CreatedAt, &newSession.ExpiresAt)
	
	if err != nil {
		log.Printf("failed to create session: %v", err)
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	cookie := &http.Cookie{
		Name: "session",
		Value: token,
		Path: "/",
		Expires: newSession.ExpiresAt,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get session cookie from the request.
		cookie, err := r.Cookie("session")
		if err != nil {
			writeJSONError(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		var authenticateSession Session

		// Look up the session using the cookie's token.
		err = db.QueryRow(
			"SELECT user_id FROM sessions WHERE token = $1 AND expires_at > $2",
			cookie.Value,
			time.Now(),
		).Scan(&authenticateSession.UserID)

		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		if err != nil {
			log.Printf("failed to query session: %v", err)
			writeJSONError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(
			r.Context(), 
			userIDKey, 
			authenticateSession.UserID,
		)

		r = r.WithContext(ctx)

		// Authentication passed, now run protected handler.
		next.ServeHTTP(w, r)
	})
}






