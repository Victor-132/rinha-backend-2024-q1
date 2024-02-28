package server

import (
	"github.com/Victor-132/rinha-backend2/service"
	"github.com/gofiber/fiber/v2"
)

type IServer interface {
	Init()
	Listen() error
}

type Server struct {
	app *fiber.App
	svc service.IService
}

func NewServer(svc service.IService) IServer {
	return &Server{svc: svc}
}

func (s *Server) Init() {
	app := fiber.New()

	app.Post("/clientes/:id/transacoes", s.svc.PostTransaction)
	app.Get("/clientes/:id/extrato", s.svc.GetStatement)

	s.app = app
}

func (s *Server) Listen() error {
	return s.app.Listen(":3000")
}
