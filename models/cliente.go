package models

type Cliente struct {
	Id     string `json:"id"`
	Saldo  int64  `json:"saldo"`
	Limite int64  `json:"limite"`
}
