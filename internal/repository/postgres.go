package repository

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type PostgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) Save(transaction *Transaction) error {
	query := `
		INSERT INTO transactions (id, valor_centavos, instante, url, is_priority)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := r.db.Exec(query, 
		transaction.ID, 
		transaction.ValorCentavos, 
		transaction.Instante, 
		transaction.URL, 
		transaction.IsPriority,
	)
	
	return err
}

func (r *PostgresTransactionRepository) GetSummary(from, to time.Time) (*PaymentsSummary, error) {
	query := `
		SELECT 
			SUM(CASE WHEN is_priority = true THEN 1 ELSE 0 END) as default_total_requests,
			SUM(CASE WHEN is_priority = true THEN valor_centavos ELSE 0 END) as default_total_amount,
			SUM(CASE WHEN is_priority = false THEN 1 ELSE 0 END) as fallback_total_requests,
			SUM(CASE WHEN is_priority = false THEN valor_centavos ELSE 0 END) as fallback_total_amount
		FROM transactions 
		WHERE instante >= $1 AND instante <= $2
	`
	
	var defaultRequests, fallbackRequests int64
	var defaultAmount, fallbackAmount int64
	
	err := r.db.QueryRow(query, from, to).Scan(
		&defaultRequests, &defaultAmount, 
		&fallbackRequests, &fallbackAmount,
	)
	if err != nil {
		return nil, err
	}
	
	summary := &PaymentsSummary{
		Default: SummaryData{
			TotalRequests: defaultRequests,
			TotalAmount:   float64(defaultAmount) / 100.0,
		},
		Fallback: SummaryData{
			TotalRequests: fallbackRequests,
			TotalAmount:   float64(fallbackAmount) / 100.0,
		},
	}
	
	return summary, nil
}