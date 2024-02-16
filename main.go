package main

import (
	"fmt"
	"net/http"

	"github.com/jpcairesf/rinha-2024-q1-go/internal/handlers"
)

func main() {
	fmt.Println("Rinha is now running...")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("POST /clientes/{id}/transacoes", handlers.PostTransacao)
	mux.HandleFunc("GET /clientes/{id}/extrato", handlers.GetExtrato)

	http.ListenAndServe(":8080", mux)
}
