package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jpcairesf/rinha-2024-q1-go/models"
	"net/http"
)

type TransacaoRequest struct {
	Valor     int64  `json:"valor" validate:"required,gt=0"`
	Tipo      string `json:"tipo" validate:"required"`
	Descricao string `json:"descricao" validate:"required,len=10"`
}

type TransacaoResponse struct {
	Saldo  int64 `json:"saldo"`
	Limite int64 `json:"limite"`
}

func PostTransacao(ctx *gin.Context) {
	var request TransacaoRequest
	err := json.NewDecoder(ctx.Request.Body).Decode(&request)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	id := ctx.Param("id")
	cliente, err := models.ExistsClienteById(id)
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

		limiteIndisponivel := cliente.Saldo < -cliente.Limite
		if limiteIndisponivel {
			http.Error(ctx.Writer, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}

	err = models.CreateTransacao(id, cliente.Saldo, request.Valor, request.Tipo, request.Descricao)

	response := TransacaoResponse{
		Saldo:  cliente.Saldo,
		Limite: cliente.Limite,
	}

	err = json.NewEncoder(ctx.Writer).Encode(response)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
