// internal/api/transaction.go
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
)

type TransactionRequest struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

func HandleTransaction(db *badger.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	urls := []string{
		"https://rinha-backend.wiremockapi.cloud/payments",
		"https://rinha-backend.wiremockapi.cloud/payments-2",
	}

	now := time.Now()
	successURL, err := TryPostToUrls(req.CorrelationID, req.Amount, now, urls)
	if err != nil {
		http.Error(w, "Payment processing failed", http.StatusInternalServerError)
		return
	}

	transactionID := uuid.New().String()
	valorCentavos := int64(req.Amount * 100)
	err = db.Update(func(txn *badger.Txn) error {
		return SaveTransaction(txn, transactionID, valorCentavos, now)
	})
	if err != nil {
		http.Error(w, "Database save failed", http.StatusInternalServerError)
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
}
