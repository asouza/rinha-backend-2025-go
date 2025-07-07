package api

import (
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger/v4"
)

func SaveTransaction(txn *badger.Txn, id string, valorCentavos int64, instante time.Time) error {
	transaction := map[string]interface{}{
		"id":       id,
		"valor":    valorCentavos,
		"instante": instante,
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	key := []byte("transaction:" + id)
	return txn.Set(key, data)
}