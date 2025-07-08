// cmd/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/asouza/rinha-backend-go/internal/api"
	"github.com/asouza/rinha-backend-go/internal/repository"
	"github.com/dgraph-io/badger/v4"
	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("./data"))
	if err != nil {
		log.Fatalf("Erro ao abrir BadgerDB: %v", err)
	}
	defer db.Close()

	transactionRepo := repository.NewBadgerTransactionRepository(db)

	urlsEnv := os.Getenv("PAYMENT_URLS")
	if urlsEnv == "" {
		log.Fatal("PAYMENT_URLS environment variable is required")
	}
	paymentURLs := strings.Split(urlsEnv, ",")

	r := chi.NewRouter()

	r.Get("/", api.HandleRoot)
	//o handler depende das urls, ficaria mais fácil de testar também.
	//Podia depender do os em si também... Podia extrapolar e criar o wrapper para expor apenas o que precisa.
	r.Post("/payments", api.HandleTransaction(transactionRepo, paymentURLs))
	r.Get("/payments-summary", api.HandlePaymentsSummary(transactionRepo))

	log.Println("Servidor rodando na porta :8080")

	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
