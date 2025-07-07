// cmd/main.go
package main

import (
	"log"
	"net/http"

	"github.com/asouza/rinha-backend-go/internal/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", api.HandleRoot)
	r.Post("/payments", api.HandleTransaction)

	log.Println("Servidor rodando na porta :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
