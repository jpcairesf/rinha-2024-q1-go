package models

import (
	"database/sql"
	"errors"
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

func TestConnection() {
	err := db.Ping()
	if err != nil {
		panic(err)
	}
}

func ExistsClienteById(id string) (*Cliente, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		// Handling transaction rollback error
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("Panic at rollback: %v", err)
		}
		// Log the panic, but allow it to propagate
		if r := recover(); r != nil {
			log.Printf("Panic at rollback but propagation is allowed: %v", r)
		}
	}()

	stmt, err := tx.Prepare("SELECT id, saldo, limite FROM cliente c WHERE id = $1")
	if err != nil {
		tx.Rollback()
		log.Printf("Panic preparing the ExistsClienteById statement: %v", err)
	}
	defer stmt.Close()

	var cliente Cliente
	err = stmt.QueryRow(id).Scan(&cliente.Id, &cliente.Saldo, &cliente.Limite)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		tx.Rollback()
		log.Printf("Panic querying the ExistsClienteById row: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Panic commiting the ExistsClienteById transaction: %v", err)
	}

	return &cliente, nil
}

func CreateTransacao(id string, saldo int64, valor int64, tipo string, descricao string) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		// Handling transaction rollback error
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("Panic at rollback: %v", err)
		}
		// Log the panic, but allow it to propagate
		if r := recover(); r != nil {
			log.Printf("Panic at rollback but propagation is allowed: %v", r)
		}
	}()

	stmtTransacao, err := tx.Prepare(
		`INSERT INTO transacao(cliente_id, valor, tipo, descricao)` +
			` VALUES ($1, $2, $3, $4)`)
	if err != nil {
		tx.Rollback()
		log.Print(err)
	}
	defer stmtTransacao.Close()

	stmtCliente, err := tx.Prepare("UPDATE cliente SET saldo = $1 WHERE cliente_id = $2")
	if err != nil {
		tx.Rollback()
		log.Print(err)
	}
	defer stmtCliente.Close()

	_, err = stmtTransacao.Exec(id, valor, tipo, descricao)
	if err != nil {
		tx.Rollback()
		log.Print(err)
	}

	_, err = stmtCliente.Exec(saldo, id)
	if err != nil {
		tx.Rollback()
		log.Print(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Panic commiting the transaction: %v", err)
	}

	return nil
}

func GetTop10TransacaoOrderByRealizadaEm(id string) ([]Transacao, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		// Handling transaction rollback error
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("Panic: %v", err)
		}
		// Log the panic, but allow it to propagate
		if r := recover(); r != nil {
			log.Printf("Panic: %v", r)
		}
	}()

	stmt, err := tx.Prepare(
		`SELECT id, cliente_id, valor, tipo, descricao` +
			` FROM transacao` +
			` WHERE cliente_id = $1` +
			` ORDER BY realizada_em DESC` +
			` LIMIT 10`)
	if err != nil {
		tx.Rollback()
		log.Print(err)
	}
	defer stmt.Close()

	var transacoes []Transacao
	rows, err := stmt.Query(id)
	for rows.Next() {
		var transacao Transacao
		if err := rows.Scan(&transacao.Id, &transacao.ClienteId,
			&transacao.Valor, &transacao.Tipo, &transacao.Descricao); err != nil {
			tx.Rollback()
			log.Print(err)
		}
		transacoes = append(transacoes, transacao)
	}
	if err = rows.Err(); err != nil {
		return transacoes, err
	}

	err = tx.Commit()
	if err != nil {
		log.Print(err)
	}

	return transacoes, nil
}
