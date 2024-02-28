package main

import (
	"github.com/Victor-132/rinha-backend2/database"
	"github.com/Victor-132/rinha-backend2/database/repository"
	"github.com/Victor-132/rinha-backend2/server"
	"github.com/Victor-132/rinha-backend2/service"
	"github.com/sirupsen/logrus"
)

func main() {
	var err error

	db := database.NewDatabase()

	err = db.Connect()
	if err != nil {
		logrus.Fatal(err)
	}

	defer db.Disconnect()

	err = db.InitDB()
	if err != nil {
		logrus.Fatal(err)
	}

	cr := repository.NewClientRepository(db)
	tr := repository.NewTransactionRepository(db)

	svc := service.NewService(cr, tr)

	srv := server.NewServer(svc)

	srv.Init()

	err = srv.Listen()
	if err != nil {
		logrus.Fatal(err)
	}
}
