package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jpcairesf/rinha-2024-q1-go/internal/db"
)

type ExtratoResponse struct {
	ExtratoSaldo      ExtratoSaldoResponse       `json:"saldo"`
	UltimasTransacoes []ExtratoTransacaoResponse `json:"ultimas_transacoes"`
}

type ExtratoSaldoResponse struct {
	Total       int64     `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int64     `json:"limite"`
}

type ExtratoTransacaoResponse struct {
	Valor       int64     `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

func GetExtrato(w http.ResponseWriter, r *http.Request) {
	db.TestConnection()

	id := r.PathValue("id")
	cliente, err := db.GetClienteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ExtratoResponse{
		ExtratoSaldo: ExtratoSaldoResponse{
			Total:       cliente.Saldo,
			DataExtrato: time.Now(),
			Limite:      cliente.Limite,
		},
		UltimasTransacoes: []ExtratoTransacaoResponse{
			{
				Valor:       100,
				Tipo:        "credito",
				Descricao:   "Depósito",
				RealizadaEm: time.Now(),
			},
			{
				Valor:       100,
				Tipo:        "debito",
				Descricao:   "Saque",
				RealizadaEm: time.Now(),
			},
		},
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
