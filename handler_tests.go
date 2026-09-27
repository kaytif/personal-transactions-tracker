package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


// Transaction creation succeeds
func TestExpenseHandler_MethodNotAllowed(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{"PATCH", http.MethodPatch},
		{"OPTIONS", http.MethodOptions},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			
			// create request
			req := httptest.NewRequest(test.method, "/transaction", nil)
			
			// create fake response writer
			rec := httptest.NewRecorder()

			// call handler
			transactionPostHandler(rec, req)

			// check result
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected 405, got %d", rec.Code)
			}
		})
	}
}






// Transfer succeeds

// Transfer atomicity

// Account balance atomicity

// Unauthenticated request

// Transacton validation

// Soft deletion



