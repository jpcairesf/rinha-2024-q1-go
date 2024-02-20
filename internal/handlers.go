package internal

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jpcairesf/rinha-2024-q1-go/api"
	"github.com/jpcairesf/rinha-2024-q1-go/internal/domain"
	"net/http"
	"time"
)

func PostTransacao(ctx *gin.Context) {
	var request api.TransacaoRequest
	err := json.NewDecoder(ctx.Request.Body).Decode(&request)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	id := ctx.Param("id")
	cliente, err := ExistsClienteById(id)
	if cliente == nil {
		http.Error(ctx.Writer, "Cliente não encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if request.Tipo == "c" {
		cliente.Saldo += request.Valor
	} else {
		cliente.Saldo -= request.Valor

		if domain.InvalidTransaction(cliente, request.Valor) {
			http.Error(ctx.Writer, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}

	err = CreateTransacao(id, cliente.Saldo, request.Valor, request.Tipo, request.Descricao)

	response := api.TransacaoResponse{
		Saldo:  cliente.Saldo,
		Limite: cliente.Limite,
	}

	err = json.NewEncoder(ctx.Writer).Encode(response)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetExtrato(ctx *gin.Context) {
	id := ctx.Param("id")
	cliente, err := ExistsClienteById(id)
	if cliente == nil {
		http.Error(ctx.Writer, "Cliente não encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	transacoes, err := GetTop10TransacaoOrderByRealizadaEm(id)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := api.ExtratoResponse{
		ExtratoSaldo: api.ExtratoSaldoResponse{
			Total:       cliente.Saldo,
			DataExtrato: time.Now(),
			Limite:      cliente.Limite,
		},
		UltimasTransacoes: []api.ExtratoTransacaoResponse{},
	}
	for _, transacao := range transacoes {
		response.UltimasTransacoes = append(response.UltimasTransacoes,
			api.ExtratoTransacaoResponse{
				Valor:       transacao.Valor,
				Tipo:        transacao.Tipo,
				Descricao:   transacao.Descricao,
				RealizadaEm: transacao.RealizadaEm,
			})
	}

	err = json.NewEncoder(ctx.Writer).Encode(response)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
