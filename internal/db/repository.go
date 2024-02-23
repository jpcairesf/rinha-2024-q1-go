package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	host     = "localhost"
	user     = "admin"
	password = "admin"
	dbname   = "rinha"
)

var (
	db                    *pgxpool.Pool
	ErrLimiteInsuficiente = errors.New("Limite insuficiente")
)

type Cliente struct {
	Id     uint8 `json:"id"`
	Saldo  int64 `json:"saldo"`
	Limite int64 `json:"limite"`
}

type Transacao struct {
	Id          uint8     `json:"id"`
	ClienteId   uint8     `json:"cliente_id"`
	Valor       int64     `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type Extrato struct {
	Saldo             SaldoExtrato       `json:"saldo"`
	UltimasTransacoes []TransacaoExtrato `json:"ultimas_transacoes"`
}

type SaldoExtrato struct {
	Total       int64     `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int64     `json:"limite"`
}

type TransacaoExtrato struct {
	Valor       int64     `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

func init() {
	conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		panic(err)
	}

	db, err = pgxpool.NewWithConfig(context.Background(), cfg)
	testConnection()
}

func testConnection() {
	err := db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("DB is now connected")
}

func CreateTransacao(ctx context.Context, transacao *Transacao) (Cliente, error) {
	var cliente Cliente
	tx, err := db.Begin(ctx)
	if err != nil {
		return cliente, err
	}
	defer func() {
		// Handling transaction rollback error
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("Panic at rollback: %v", err)
		}
		// Log the panic, but allow it to propagate
		if r := recover(); r != nil {
			log.Printf("Panic at rollback but propagation is allowed: %v", r)
		}
	}()

	// if tipo == 'd' AND saldo + limite < valor
	err = tx.QueryRow(ctx,
		"SELECT limite, saldo FROM cliente WHERE id = $1 FOR UPDATE", transacao.ClienteId,
	).Scan(&cliente.Limite, &cliente.Saldo)
	if err != nil {
		return cliente, err
	}
	fmt.Println(cliente)
	if transacao.Tipo == "c" {
		cliente.Saldo += transacao.Valor
	} else {
		if cliente.Saldo+cliente.Limite < transacao.Valor {
			return cliente, ErrLimiteInsuficiente
		}
		cliente.Saldo -= transacao.Valor
	}

	batch := &pgx.Batch{}
	batch.Queue("UPDATE cliente SET saldo = $1 WHERE id = $2", cliente.Saldo, transacao.ClienteId)
	batch.Queue(`INSERT INTO transacao (cliente_id, valor, tipo, descricao, realizada_em)`+
		` VALUES ($1, $2, $3, $4, $5)`,
		transacao.ClienteId, transacao.Valor, transacao.Tipo, transacao.Descricao, transacao.RealizadaEm,
	)
	fmt.Println(cliente)

	result := tx.SendBatch(ctx, batch)
	if err := result.Close(); err != nil {
		return Cliente{}, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return Cliente{}, err
	}

	return cliente, nil
}

func GetTop10TransacaoOrderByRealizadaEm(ctx context.Context, id uint8) (Extrato, error) {
	var saldo SaldoExtrato
	err := db.QueryRow(ctx,
		"SELECT limite, saldo FROM cliente WHERE id = $1", id,
	).Scan(&saldo.Limite, &saldo.Total)
	if err != nil {
		return Extrato{}, err
	}

	results, err := db.Query(ctx,
		`SELECT valor, tipo, descricao, realizada_em`+
			` FROM transacao`+
			` WHERE cliente_id = $1`+
			` ORDER BY realizada_em DESC`+
			` LIMIT 10`, id)
	if err != nil {
		return Extrato{}, err
	}

	var transacoes []TransacaoExtrato
	for results.Next() {
		var transacao TransacaoExtrato
		err = results.Scan(&transacao.Valor, &transacao.Tipo, &transacao.Descricao, &transacao.RealizadaEm)
		if err != nil {
			return Extrato{}, err
		}
		transacoes = append(transacoes, transacao)
	}

	return Extrato{saldo, transacoes}, nil
}
