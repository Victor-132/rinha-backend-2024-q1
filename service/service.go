package service

import (
	"math"
	"time"

	"github.com/Victor-132/rinha-backend2/database"
	"github.com/Victor-132/rinha-backend2/database/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IService interface {
	PostTransaction(c *fiber.Ctx) error
	GetStatement(c *fiber.Ctx) error
}

type Service struct {
	cr repository.IClientRepository
	tr repository.ITransactionRepository
}

func NewService(cr repository.IClientRepository, tr repository.ITransactionRepository) IService {
	return &Service{cr, tr}
}

func (s *Service) PostTransaction(c *fiber.Ctx) (err error) {
	var id int

	id, err = c.ParamsInt("id")
	if err != nil {
		logrus.Error(err)
		return
	}

	var cl *database.Client
	cl, err = s.cr.GetClientById(id)
	if cl == nil && err == nil {
		return c.Status(fiber.StatusNotFound).SendString("cliente não encontrado.")
	}

	if err != nil {
		logrus.Error(err)
		return
	}

	m := PostTransaction{}
	err = c.BodyParser(&m)
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	err = m.validate()
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	val := m.Value
	if m.Type == "c" {
		cl.Balance += val
	}

	if m.Type == "d" {
		cl.Balance -= val
		val *= -1

		if float64(cl.Limit)-math.Abs(float64(cl.Balance)) < 0 {
			return c.Status(fiber.StatusUnprocessableEntity).
				SendString("limite insuficiente para realizar a trasação.")
		}
	}

	err = s.cr.SetClientBalanceByClientId(id, val)
	if err != nil {
		return
	}

	t := database.Transaction{
		ClientId:    cl.Id,
		Value:       m.Value,
		Type:        m.Type,
		Description: m.Description,
		Date:        time.Now(),
	}

	err = s.tr.CreateTransaction(t)
	if err != nil {
		return
	}

	ret := ResponseTransaction{
		Limite: cl.Limit,
		Saldo:  cl.Balance,
	}

	return c.JSON(ret)
}

func (s *Service) GetStatement(c *fiber.Ctx) (err error) {
	var id int

	id, err = c.ParamsInt("id")
	if err != nil {
		logrus.Error(err)
		return
	}

	var cl *database.Client
	cl, err = s.cr.GetClientById(id)
	if cl == nil && err == nil {
		return c.Status(fiber.StatusNotFound).SendString("cliente não encontrado.")
	}

	if err != nil {
		logrus.Error(err)
		return
	}

	var t []database.Transaction
	t, err = s.tr.GetTransactionsByClientId(cl.Id)
	if err != nil {
		logrus.Error(err)
		return
	}

	st := Statement{
		Balance: Balance{
			Total: cl.Balance,
			Limit: cl.Limit,
			Date:  time.Now(),
		},
		LastTransactions: []LastTransaction{},
	}

	i := LastTransaction{}
	for _, t := range t {
		i.Date = t.Date
		i.Description = t.Description
		i.Type = t.Type
		i.Value = t.Value

		st.LastTransactions = append(st.LastTransactions, i)
	}

	return c.JSON(st)
}
