package domain

type Cliente struct {
	Id     string `json:"id"`
	Saldo  int64  `json:"saldo"`
	Limite int64  `json:"limite"`
}

func InvalidTransaction(c *Cliente, valor int64) bool {
	return c.Saldo-valor < c.Limite
}
