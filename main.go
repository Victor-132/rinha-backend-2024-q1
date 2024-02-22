package main

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/clientes/:id/transacoes", post)

	app.Get("/clientes/:id/extrato", func(c *fiber.Ctx) error {
		return nil
	})

	app.Listen(":3000")
}

type transaction string

const (
	debit  = "d"
	credit = "c"
)

type Body struct {
	Valor     int
	Tipo      transaction
	Descricao string
}

func (b *Body) validate() error {
	if b.Valor <= 0 {
		return errors.New("o valor deve ser um valor inteiro maior que 0")
	}

	if b.Tipo != debit && b.Tipo != credit {
		return errors.New("o tipo de transação deve ser c (crédito) ou d (débito)")
	}

	if len(b.Descricao) < 1 || len(b.Descricao) > 10 {
		return errors.New("a descrição deve ter entre 1 e 10 caracteres")
	}

	return nil
}

func post(c *fiber.Ctx) error {
	b := Body{}

	if err := c.BodyParser(&b); err != nil {
		return err
	}

	return b.validate()
}
