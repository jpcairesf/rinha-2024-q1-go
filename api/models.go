package api

import "time"

type TransacaoRequest struct {
	Valor     int64  `json:"valor" validate:"required,gt=0"`
	Tipo      string `json:"tipo" validate:"required"`
	Descricao string `json:"descricao" validate:"required,len=10"`
}

type TransacaoResponse struct {
	Saldo  int64 `json:"saldo"`
	Limite int64 `json:"limite"`
}

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
