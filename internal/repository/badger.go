package repository

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type BadgerTransactionRepository struct {
	db *badger.DB
}

func NewBadgerTransactionRepository(db *badger.DB) *BadgerTransactionRepository {
	return &BadgerTransactionRepository{db: db}
}

func (r *BadgerTransactionRepository) Save(transaction *Transaction) error {
	return r.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(transaction)
		if err != nil {
			return err
		}

		key := []byte("transaction:" + transaction.ID)
		return txn.Set(key, data)
	})
}

func (r *BadgerTransactionRepository) GetSummary(from, to time.Time) (*PaymentsSummary, error) {
	summary := &PaymentsSummary{}

	err := r.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte("transaction:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var transaction Transaction
				if err := json.Unmarshal(val, &transaction); err != nil {
					return err
				}

				if transaction.Instante.After(from) && transaction.Instante.Before(to) {
					amount := float64(transaction.ValorCentavos) / 100.0
					
					if transaction.IsPriority {
						summary.Default.TotalRequests++
						summary.Default.TotalAmount += amount
					} else {
						summary.Fallback.TotalRequests++
						summary.Fallback.TotalAmount += amount
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return summary, err
}
