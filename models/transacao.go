package models

import "time"

type Transacao struct {
	Id          string    `json:"id"`
	ClienteId   string    `json:"cliente_id"`
	Valor       int64     `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}
