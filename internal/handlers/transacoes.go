package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jpcairesf/rinha-2024-q1-go/internal/db"
	"net/http"
	"strconv"
	"time"
)

type TransacaoRequest struct {
	Valor     float64 `json:"valor"`
	Tipo      string  `json:"tipo"`
	Descricao string  `json:"descricao"`
}

type TransacaoResponse struct {
	Saldo  int64 `json:"saldo"`
	Limite int64 `json:"limite"`
}

func (t *TransacaoRequest) isNotValid() bool {
	tipoValido := t.Tipo == "c" || t.Tipo == "d"
	descricaoValida := len(t.Descricao) > 0 && len(t.Descricao) <= 10
	valorValido := t.Valor >= 0 && t.Valor == float64(int64(t.Valor))
	return !(tipoValido && descricaoValida && valorValido)
}

func PostTransacao(w http.ResponseWriter, r *http.Request) {
	var request TransacaoRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if request.isNotValid() {
		http.Error(w, "Validation failed", http.StatusUnprocessableEntity)
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transacao := db.Transacao{
		ClienteId:   uint8(id),
		Valor:       int64(request.Valor),
		Tipo:        request.Tipo,
		Descricao:   request.Descricao,
		RealizadaEm: time.Now(),
	}

	cliente, err := db.CreateTransacao(context.Background(), &transacao)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, db.ErrLimiteInsuficiente) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}

	response := TransacaoResponse{Saldo: cliente.Saldo, Limite: cliente.Limite}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
