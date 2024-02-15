package handlers

import (
	"encoding/json"
	"fmt"
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
	// Decode the request body into the payload struct
	var payload TransacaoRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	id := r.PathValue("id")
	cliente, err := db.GetClienteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(cliente)

	response := TransacaoResponse{
		Saldo:  cliente.Saldo,
		Limite: cliente.Limite,
	}

	// response := TransacaoResponse{
	// 	Saldo:  1000,
	// 	Limite: 1000,
	// }

	// Encode the response payload and send it back
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
