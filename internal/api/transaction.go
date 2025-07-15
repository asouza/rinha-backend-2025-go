// internal/api/transaction.go
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/asouza/rinha-backend-go/internal/repository"
	"github.com/google/uuid"
)

type TransactionRequest struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

func HandleTransaction(repo repository.TransactionRepository, externalClient ExternalServiceClient, paymentURLs []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TransactionRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if len(paymentURLs) == 0 {
			log.Println("Payment URLs not configured")
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}
		urls := paymentURLs

		now := time.Now().UTC()
		successURL, err := externalClient.TryPostToUrls(req.CorrelationID, req.Amount, now, urls)
		if err != nil {
			http.Error(w, "Payment processing failed", http.StatusInternalServerError)
			return
		}

		transactionID := uuid.New().String()
		valorCentavos := int64(req.Amount * 100)

		isPriority := successURL == urls[0]

		transaction := &repository.Transaction{
			ID:            transactionID,
			ValorCentavos: valorCentavos,
			Instante:      now,
			URL:           successURL,
			IsPriority:    isPriority,
		}

		err = repo.Save(transaction)
		if err != nil {
			http.Error(w, "Database save failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func HandlePaymentsSummary(repo repository.TransactionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")

		var from, to time.Time
		var err error

		if fromStr == "" {
			from = time.Now().AddDate(-20, 0, 0)
		} else {
			from, err = time.Parse("2006-01-02T15:04:05.000Z", fromStr)
			if err != nil {
				http.Error(w, "Invalid 'from' date format. Use ISO 8601 format", http.StatusBadRequest)
				return
			}
		}

		if toStr == "" {
			to = time.Now().AddDate(20, 0, 0)
		} else {
			to, err = time.Parse("2006-01-02T15:04:05.000Z", toStr)
			if err != nil {
				http.Error(w, "Invalid 'to' date format. Use ISO 8601 format", http.StatusBadRequest)
				return
			}
		}

		if from.After(to) {
			http.Error(w, "'from' date must be before 'to' date", http.StatusBadRequest)
			return
		}

		//aqui tem aquele acoplamento que não queremos.
		//o retorno do repository já está sendo enviado para o client.. tinha que transformar.
		//mas se transformar, tem que iterar de novo, e agora?


		summary, err := repo.GetSummary(from, to)
		if err != nil {
			http.Error(w, "Failed to get summary", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(summary)
	}
}
