package main

func getBalance(id int) (float64, error) {
	var balance float64
	err := db.QueryRow(
		"SELECT COALESC(SUM(amount), 0) FROM transactions WHERE account_id = $1 AND deleted_at is NULL",
		id, 
	).Scan(&balance)

	return balance, err
}