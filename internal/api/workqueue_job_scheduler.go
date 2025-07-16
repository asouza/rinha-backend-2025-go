package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"k8s.io/client-go/util/workqueue"
)

// WorkqueueJobScheduler implements JobScheduler using Kubernetes workqueue
type WorkqueueJobScheduler struct {
	queue       workqueue.TypedDelayingInterface[RetryPaymentJob]
	httpClient  *http.Client
	baseURL     string
	workers     int
	stopCh      chan struct{}
	wg          sync.WaitGroup
	maxRetries  int
	initialDelay time.Duration
}

// NewWorkqueueJobScheduler creates a new workqueue-based job scheduler
func NewWorkqueueJobScheduler(baseURL string, workers int, maxRetries int, initialDelay time.Duration) *WorkqueueJobScheduler {
	return &WorkqueueJobScheduler{
		queue:       workqueue.NewTypedDelayingQueue[RetryPaymentJob](),
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		baseURL:     baseURL,
		workers:     workers,
		stopCh:      make(chan struct{}),
		maxRetries:  maxRetries,
		initialDelay: initialDelay,
	}
}

// ScheduleRetryPayment schedules a retry payment job
func (w *WorkqueueJobScheduler) ScheduleRetryPayment(payload TransactionRequest, delay time.Duration) error {
	job := RetryPaymentJob{
		Payload:      payload,
		AttemptCount: 0,
		MaxRetries:   w.maxRetries,
	}

	w.queue.AddAfter(job, delay)
	log.Printf("Scheduled retry payment job for correlation ID: %s, delay: %v", payload.CorrelationID, delay)
	return nil
}

// Start starts the job scheduler workers
func (w *WorkqueueJobScheduler) Start() error {
	log.Printf("Starting %d workers for retry payment jobs", w.workers)
	
	for i := 0; i < w.workers; i++ {
		w.wg.Add(1)
		go w.worker()
	}
	
	return nil
}

// Stop stops the job scheduler
func (w *WorkqueueJobScheduler) Stop() {
	log.Println("Stopping job scheduler...")
	close(w.stopCh)
	w.queue.ShutDown()
	w.wg.Wait()
	log.Println("Job scheduler stopped")
}

// worker processes jobs from the queue
func (w *WorkqueueJobScheduler) worker() {
	defer w.wg.Done()
	
	for {
		select {
		case <-w.stopCh:
			return
		default:
			item, shutdown := w.queue.Get()
			if shutdown {
				return
			}
			
			w.processJob(item)
			w.queue.Done(item)
		}
	}
}

// processJob processes a single retry payment job
func (w *WorkqueueJobScheduler) processJob(job RetryPaymentJob) {
	
	job.AttemptCount++
	log.Printf("Processing retry payment job (attempt %d/%d) for correlation ID: %s", 
		job.AttemptCount, job.MaxRetries, job.Payload.CorrelationID)
	
	err := w.sendPaymentRequest(job.Payload)
	if err != nil {
		log.Printf("Retry payment failed for correlation ID %s: %v", job.Payload.CorrelationID, err)
		
		if job.AttemptCount < job.MaxRetries {
			// Calculate exponential backoff delay
			delay := w.calculateBackoffDelay(job.AttemptCount)
			w.queue.AddAfter(job, delay)
			log.Printf("Scheduled retry %d/%d for correlation ID %s with delay %v", 
				job.AttemptCount+1, job.MaxRetries, job.Payload.CorrelationID, delay)
		} else {
			log.Printf("Max retries reached for correlation ID %s", job.Payload.CorrelationID)
		}
	} else {
		log.Printf("Retry payment successful for correlation ID: %s", job.Payload.CorrelationID)
	}
}

// sendPaymentRequest sends the payment request to the /payments endpoint
func (w *WorkqueueJobScheduler) sendPaymentRequest(payload TransactionRequest) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	url := fmt.Sprintf("%s/payments", w.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := w.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}
	
	return nil
}

// calculateBackoffDelay calculates exponential backoff delay
func (w *WorkqueueJobScheduler) calculateBackoffDelay(attemptCount int) time.Duration {
	backoffFactor := time.Duration(1 << (attemptCount - 1)) // 2^(attemptCount-1)
	delay := w.initialDelay * backoffFactor
	
	// Cap the delay at 5 minutes
	maxDelay := 5 * time.Minute
	if delay > maxDelay {
		delay = maxDelay
	}
	
	return delay
}