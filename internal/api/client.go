// internal/api/client.go
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type PostRequest struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
	RequestedAt   string  `json:"requestedAt"`
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

	for i, url := range urls {
		log.Printf("Tentando conectar na URL %d: %s", i+1, url)
		
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Erro ao conectar na URL %s: %v", url, err)
			continue
		}

		log.Printf("Resposta da URL %s: Status %d", url, resp.StatusCode)
		
		// Ler corpo da resposta para debug
		if resp.Body != nil {
			body, readErr := io.ReadAll(resp.Body)
			if readErr == nil {
				log.Printf("Corpo da resposta de %s: %s", url, string(body))
			}
			resp.Body.Close()
		}

		if resp.StatusCode == http.StatusOK {
			log.Printf("Sucesso na URL: %s", url)
			return url, nil
		} else {
			log.Printf("URL %s falhou com status: %d", url, resp.StatusCode)
		}
	}

	log.Printf("Todas as URLs falharam. Total de URLs testadas: %d", len(urls))
	return "", errors.New("all URLs failed")
}
