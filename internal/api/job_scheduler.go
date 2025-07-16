package api

import (
	"time"
)

// JobScheduler defines the interface for scheduling retry jobs
type JobScheduler interface {
	ScheduleRetryPayment(payload TransactionRequest, delay time.Duration) error
	Start() error
	Stop()
}

// RetryPaymentJob represents a job to retry a payment with metadata
type RetryPaymentJob struct {
	Payload      TransactionRequest `json:"payload"`
	AttemptCount int                `json:"attemptCount"`
	MaxRetries   int                `json:"maxRetries"`
}