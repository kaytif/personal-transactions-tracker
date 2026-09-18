package main

import (
	"encoding/json"
	"net/http"

)


func transactionPostHandler(w http.ResponseWriter, r *http.Request) {
	// post is adding so no id check id will be generated automatically
	// name needs to be checked against the things
	// okay if transaction name is unique
	// you needd thee  account id that links the transaction to your 
	