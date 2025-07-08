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

type TransactionRepository interface {
	Save(transaction *Transaction) error
}
