package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jpcairesf/rinha-2024-q1-go/internal/db"
)

type TransacaoRequest struct {
	Valor     int64  `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type TransacaoResponse struct {
	Saldo  int64 `json:"saldo"`
	Limite int64 `json:"limite"`
}

func PostTransacao(w http.ResponseWriter, r *http.Request) {
	var request TransacaoRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	id := r.PathValue("id")
	cliente, err := db.ExistsClienteById(id)
	if cliente == nil {
		http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if request.Tipo == "c" {
		cliente.Saldo += request.Valor
	} else {
		cliente.Saldo -= request.Valor

		limiteIndisponivel := cliente.Saldo < -cliente.Limite
		if limiteIndisponivel {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}

	err = db.CreateTransacao(id, cliente.Saldo, request.Valor, request.Tipo, request.Descricao)

	response := TransacaoResponse{
		Saldo:  cliente.Saldo,
		Limite: cliente.Limite,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
