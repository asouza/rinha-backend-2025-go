// internal/api/client.go
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type PostRequest struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
	RequestedAt   string  `json:"requestedAt"`
}

type ExternalServiceClient interface {
	TryPostToUrls(correlationID string, amount float64, requestedAt time.Time, urls []string) (string, error)
}

type HTTPClient struct {
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

func (c *HTTPClient) TryPostToUrls(correlationID string, amount float64, requestedAt time.Time, urls []string) (string, error) {
	payload := PostRequest{
		CorrelationID: correlationID,
		Amount:        amount,
		RequestedAt:   requestedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	for i, url := range urls {
		var timeout time.Duration
		if i == 0 {
			timeout = 500 * time.Millisecond
		} else {
			timeout = 5000 * time.Millisecond
		}

		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Resposta e erro na comunicacao com url %v: %v %v", url, resp, err)
			continue
		}

		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return url, nil
		}
	}

	return "", errors.New("all URLs failed")
}
