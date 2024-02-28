package service

import (
	"fmt"
	"time"
)

const (
	debit  = "d"
	credit = "c"
)

type PostTransaction struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

func (b *PostTransaction) validate() error {
	if b.Value <= 0 {
		return fmt.Errorf(
			"o valor deve ser um valor inteiro maior que 0 - valor: %d",
			b.Value,
		)
	}

	if b.Type != debit && b.Type != credit {
		return fmt.Errorf(
			"o tipo de transação deve ser c (crédito) ou d (débito) - tipo: %s",
			b.Description,
		)
	}

	if len(b.Description) < 1 || len(b.Description) > 10 {
		return fmt.Errorf(
			"a descrição deve ter entre 1 e 10 caracteres - descricao: %s (%d)",
			b.Description,
			len(b.Description),
		)
	}

	return nil
}

type ResponseTransaction struct {
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
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
