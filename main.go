package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jpcairesf/rinha-2024-q1-go/internal/handlers"
)

func main() {
	fmt.Println("Rinha is now running...")
	mux := http.NewServeMux()
	mux.HandleFunc("POST /clientes/{id}/transacoes", handlers.PostTransacao)
	mux.HandleFunc("GET /clientes/{id}/extrato", handlers.GetExtrato)

	err := http.ListenAndServe(":"+os.Getenv("RINHA_PORT"), mux)
	if err != nil {
		return
	}
}
