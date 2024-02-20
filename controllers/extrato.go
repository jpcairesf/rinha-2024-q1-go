package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jpcairesf/rinha-2024-q1-go/models"
	"net/http"
	"time"
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

func GetExtrato(c *gin.Context) {
	id := c.Param("id")
	cliente, err := models.ExistsClienteById(id)
	if cliente == nil {
		http.Error(c.Writer, "Cliente não encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	transacoes, err := models.GetTop10TransacaoOrderByRealizadaEm(id)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ExtratoResponse{
		ExtratoSaldo: ExtratoSaldoResponse{
			Total:       cliente.Saldo,
			DataExtrato: time.Now(),
			Limite:      cliente.Limite,
		},
		UltimasTransacoes: []ExtratoTransacaoResponse{},
	}
	for _, transacao := range transacoes {
		response.UltimasTransacoes = append(response.UltimasTransacoes,
			ExtratoTransacaoResponse{
				Valor:       transacao.Valor,
				Tipo:        transacao.Tipo,
				Descricao:   transacao.Descricao,
				RealizadaEm: transacao.RealizadaEm,
			})
	}

	err = json.NewEncoder(c.Writer).Encode(response)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
