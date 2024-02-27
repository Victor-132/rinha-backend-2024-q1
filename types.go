package main

import (
	"fmt"
	"time"
)

type Client struct {
	Id      int
	Limit   int
	Balance int
}

type Transaction struct {
	ClientId    int
	Value       int
	Type        string
	Description string
	Date        time.Time
}

type Statement struct {
	Balance          Balance           `json:"saldo"`
	LastTransactions []LastTransaction `json:"ultimas_transacoes"`
}

type Balance struct {
	Total int       `json:"total"`
	Limit int       `json:"limite"`
	Date  time.Time `json:"data_extrato"`
}

type LastTransaction struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	Date        time.Time `json:"realizada_em"`
}

const (
	debit  = "d"
	credit = "c"
)

type Body struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

func (b *Body) validate() error {
	if b.Value <= 0 {
		return fmt.Errorf("o valor deve ser um valor inteiro maior que 0 - valor: %d", b.Value)
	}

	if b.Type != debit && b.Type != credit {
		return fmt.Errorf("o tipo de transação deve ser c (crédito) ou d (débito) - tipo: %s", b.Description)
	}

	if len(b.Description) < 1 || len(b.Description) > 10 {
		return fmt.Errorf("a descrição deve ter entre 1 e 10 caracteres - descricao: %s (%d)", b.Description, len(b.Description))
	}

	return nil
}
