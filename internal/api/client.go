// internal/api/client.go
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type PostRequest struct {
	CorrelationID string `json:"correlationId"`
	Amount        float64 `json:"amount"`
	RequestedAt   string `json:"requestedAt"`
}

func TryPostToUrls(correlationID string, amount float64, requestedAt time.Time, urls []string) (string, error) {
	payload := PostRequest{
		CorrelationID: correlationID,
		Amount:        amount,
		RequestedAt:   requestedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	
	fmt.Printf("Payload para endpoint externo: %s\n", string(jsonData))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for _, url := range urls {
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}

		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return url, nil
		}
	}

	return "", errors.New("all URLs failed")
}