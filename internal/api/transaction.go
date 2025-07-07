// internal/api/transaction.go
package api

import (
	"encoding/json"
	"net/http"
)

type TransactionRequest struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

func HandleTransaction(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"correlationId": req.CorrelationID,
		"amount":        req.Amount,
		"status":        "received",
	}
	
	json.NewEncoder(w).Encode(response)
}