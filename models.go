package main

import "time"

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
	CategoryID int `json:"category_id"`
	DeletedAt *time.Time `json:"deleted_at"`
}