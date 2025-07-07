package repository

import (
	"encoding/json"

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
		transactionMap := map[string]interface{}{
			"id":       transaction.ID,
			"valor":    transaction.ValorCentavos,
			"instante": transaction.Instante,
		}

		data, err := json.Marshal(transactionMap)
		if err != nil {
			return err
		}

		key := []byte("transaction:" + transaction.ID)
		return txn.Set(key, data)
	})
}
