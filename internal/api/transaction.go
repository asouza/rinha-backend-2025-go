// internal/api/transaction.go
package api

import (
	"encoding/json"
	"net/http"
	"time"
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

	urls := []string{
		"https://rinha-backend.wiremockapi.cloud/payments",
		"https://rinha-backend.wiremockapi.cloud/payments-2",
	}

	successURL, err := TryPostToUrls(req.CorrelationID, req.Amount, time.Now(), urls)
	if err != nil {
		http.Error(w, "Payment processing failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"correlationId": req.CorrelationID,
		"amount":        req.Amount,
		"status":        "processed",
		"processedBy":   successURL,
	}

	json.NewEncoder(w).Encode(response)
}
