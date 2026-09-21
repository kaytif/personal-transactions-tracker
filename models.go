package main

import (
	"time"
	"vendor/golang.org/x/net/idna"
)

// Expense represents one expense

// We use json: because while converting to and from JSOn
// it tells GO what tag to store or return information in

// We store only name and id because having the account balance
// can create synchornization problems so it is better to calculate as when needed
type Account struct {
	ID int `json:"id"`
	Name string `json:"name"`
	DeletedAt *time.Time `json:"deleted_at"`
	Balance float64 `json:"balance"`
}

type Category struct {
	ID int `json:"id"`
	Name string `json:"name"`
}


type Transaction struct {
	ID int `json:"id"`
	AccountID int `json:"account_id"`
	Name string `json:"name"`
	Amount float64 `json:"amount"`
	Date string `json:"transaction_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	// A transaction may or may not have a category. But an actual category always has an ID.
	CategoryID *int `json:"category_id"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// this is what is received temporarily
// we hash the password but never store the password itself
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


// this is what our system stores. our system must always store the hash 
// of the pasword and the not the actual password itself for security reasons
type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	PasswordHash string `json:"-"` // Never send hash in JSON responses
}

