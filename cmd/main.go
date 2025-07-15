// cmd/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/asouza/rinha-backend-go/internal/api"
	"github.com/asouza/rinha-backend-go/internal/repository"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	// Get database connection parameters from environment variables
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER environment variable is required")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD environment variable is required")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is required")
	}

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão com PostgreSQL: %v", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(2 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao conectar com PostgreSQL: %v", err)
	}

	transactionRepo := repository.NewPostgresTransactionRepository(db)
	externalClient := api.NewHTTPClient()

	urlsEnv := os.Getenv("PAYMENT_URLS")
	if urlsEnv == "" {
		log.Fatal("PAYMENT_URLS environment variable is required")
	}
	paymentURLs := strings.Split(urlsEnv, ",")

	r := chi.NewRouter()

	r.Get("/", api.HandleRoot)
	//o handler depende das urls, ficaria mais fácil de testar também.
	//Podia depender do os em si também... Podia extrapolar e criar o wrapper para expor apenas o que precisa.
	r.Post("/payments", api.HandleTransaction(transactionRepo, externalClient, paymentURLs))
	r.Get("/payments-summary", api.HandlePaymentsSummary(transactionRepo))

	log.Println("Servidor rodando na porta :9999")

	err = http.ListenAndServe(":9999", r)

	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
