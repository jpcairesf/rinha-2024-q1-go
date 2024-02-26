package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jpcairesf/rinha-2024-q1-go/internal/db"
	"log"
	"net/http"
	"strconv"
)

func GetExtrato(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	extrato, err := db.GetTop10TransacaoOrderByRealizadaEm(context.Background(), uint8(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Error retrieving extrato: %v", err)
	}

	err = json.NewEncoder(w).Encode(extrato)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
