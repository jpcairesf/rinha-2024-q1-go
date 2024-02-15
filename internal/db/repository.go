package db

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "admin"
	dbname   = "rinha"
)

var db *sql.DB

func init() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

func TestConnection() {
	err := db.Ping()
	if err != nil {
		panic(err)
	}
}

type Cliente struct {
	Id     string `json:"id"`
	Saldo  int64  `json:"saldo"`
	Limite int64  `json:"limite"`
}

func GetClienteByID(id string) (*Cliente, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the statement
	stmt, err := tx.Prepare("SELECT c.cliente_id, c.saldo, c.limite FROM cliente c WHERE cliente_id = $1")
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmt.Close()

	var cliente Cliente
	err = stmt.QueryRow(id).Scan(&cliente.Saldo, &cliente.Limite)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		tx.Rollback()
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return &cliente, nil
}
