package repository

import (
	"time"
)

type Transaction struct {
	ID            string    `json:"id"`
	ValorCentavos int64     `json:"valor"`
	Instante      time.Time `json:"instante"`
	URL           string    `json:"url"`
	IsPriority    bool      `json:"isPriority"`
}

type SummaryData struct {
	TotalRequests int64   `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

type PaymentsSummary struct {
	Default  SummaryData `json:"default"`
	Fallback SummaryData `json:"fallback"`
}

type TransactionRepository interface {
	Save(transaction *Transaction) error
	GetSummary(from, to time.Time) (*PaymentsSummary, error)
}
