package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Transaction rejects invalid methods
func TestHandlerTransaction_MethodNotAllowed(t *testing.T) {

	// create test cases
	tests := []struct {
		name   string
		method string
	}{
		{"PATCH", http.MethodPatch},
		{"OPTIONS", http.MethodOptions},
	}

	// run each test case
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// create request
			req := httptest.NewRequest(test.method, "/transactions", nil)

			// create fake response writer
			rec := httptest.NewRecorder()

			// call transaction handler
			handlerTransaction(rec, req)

			// check response is 405
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405, got %d", rec.Code)
			}
		})
	}
}

// Account rejects invalid methods
func TestHandlerAccount_MethodNotAllowed(t *testing.T) {

	// create test cases
	tests := []struct {
		name   string
		method string
	}{
		{"PATCH", http.MethodPatch},
		{"OPTIONS", http.MethodOptions},
	}

	// run each test case
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// create request
			req := httptest.NewRequest(test.method, "/accounts", nil)

			// create fake response writer
			rec := httptest.NewRecorder()

			// call account handler
			handlerAccount(rec, req)

			// check response is 405
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405, got %d", rec.Code)
			}
		})
	}
}

// Category rejects invalid methods
func TestHandlerCategory_MethodNotAllowed(t *testing.T) {

	// create test cases
	tests := []struct {
		name   string
		method string
	}{
		{"PATCH", http.MethodPatch},
		{"OPTIONS", http.MethodOptions},
	}

	// run each test case
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// create request
			req := httptest.NewRequest(test.method, "/categories", nil)

			// create fake response writer
			rec := httptest.NewRecorder()

			// call category handler
			handlerCategory(rec, req)

			// check response is 405
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405, got %d", rec.Code)
			}
		})
	}
}

// Transfer rejects invalid methods
func TestHandlerTransfer_MethodNotAllowed(t *testing.T) {

	// create test cases
	tests := []struct {
		name   string
		method string
	}{
		{"PATCH", http.MethodPatch},
		{"OPTIONS", http.MethodOptions},
	}

	// run each test case
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// create request
			req := httptest.NewRequest(test.method, "/transfers", nil)

			// create fake response writer
			rec := httptest.NewRecorder()

			// call transfer handler
			handlerTransfer(rec, req)

			// check response is 405
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405, got %d", rec.Code)
			}
		})
	}
}
