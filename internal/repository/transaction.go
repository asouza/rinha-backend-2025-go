package repository

import (
	"time"
)

type Transaction struct {
	ID            string    `json:"id"`
	ValorCentavos int64     `json:"valor"`
	Instante      time.Time `json:"instante"`
}

type TransactionRepository interface {
	Save(transaction *Transaction) error
}
