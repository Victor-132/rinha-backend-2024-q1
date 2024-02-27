package main

import (
	"context"
	"errors"
	"math"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func main() {
	clientOptions := options.Client().
		ApplyURI(os.Getenv("DATABASE_URL"))

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	db = client.Database("rinha")

	cc := db.Collection("clients")
	tc := db.Collection("transactions")

	cc.DeleteMany(context.TODO(), bson.M{}, nil)
	tc.DeleteMany(context.TODO(), bson.M{}, nil)

	clients := []Client{
		{Id: 1, Limit: 100000, Balance: 0},
		{Id: 2, Limit: 80000, Balance: 0},
		{Id: 3, Limit: 1000000, Balance: 0},
		{Id: 4, Limit: 10000000, Balance: 0},
		{Id: 5, Limit: 500000, Balance: 0},
	}

	filter := bson.M{}
	update := bson.M{}

	for _, c := range clients {
		filter["id"] = c.Id
		update["$set"] = c

		opt := options.UpdateOptions{}
		opt.SetUpsert(true)

		cc.UpdateOne(context.TODO(), filter, update, &opt)
	}

	cc.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "id", Value: 1}}})
	tc.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "clientid", Value: 1}, {Key: "date", Value: 1}}})

	app := fiber.New()

	app.Post("/clientes/:id/transacoes", post)

	app.Get("/clientes/:id/extrato", get)

	app.Listen(":3000")
}

func post(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	cc := db.Collection("clients")

	filter := bson.M{
		"id": id,
	}

	res := cc.FindOne(context.TODO(), filter, nil)

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).SendString("cliente não encontrado.")
		}

		logrus.Error(res.Err())
		return res.Err()
	}

	client := Client{}

	if err := res.Decode(&client); err != nil {
		logrus.Error(err)
		return err
	}

	if err != nil {
		logrus.Error(err)
		return err
	}

	b := Body{}

	if err := c.BodyParser(&b); err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	if err := b.validate(); err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	amount := b.Value
	if b.Type == "c" {
		client.Balance += amount
	}

	if b.Type == "d" {
		client.Balance -= amount
		amount *= -1

		if float64(client.Limit)-math.Abs(float64(client.Balance)) < 0 {
			return c.Status(fiber.StatusUnprocessableEntity).SendString("limite insuficiente para realizar a trasação.")
		}
	}

	update := bson.M{
		"$inc": bson.M{
			"balance": amount,
		},
	}

	resUpdate, errUpdate := cc.UpdateOne(context.TODO(), filter, update, nil)

	if errUpdate != nil {
		logrus.Error(errUpdate)
		return errUpdate
	}

	if resUpdate.ModifiedCount == 0 {
		err := errors.New("nenhum saldo foi atualizado")
		logrus.Error(err)
		return err
	}

	t := Transaction{
		ClientId:    client.Id,
		Value:       b.Value,
		Type:        b.Type,
		Description: b.Description,
		Date:        time.Now(),
	}

	tc := db.Collection("transactions")

	_, errT := tc.InsertOne(context.TODO(), t, nil)

	if errT != nil {
		logrus.Error(errT)
		return errT
	}

	ret := struct {
		Limite int `json:"limite"`
		Saldo  int `json:"saldo"`
	}{
		Limite: client.Limit,
		Saldo:  client.Balance,
	}

	return c.JSON(ret)
}

func get(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		logrus.Error(err)
		return err
	}

	cc := db.Collection("clients")

	filter := bson.M{
		"id": id,
	}

	res := cc.FindOne(context.TODO(), filter, nil)

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).SendString("cliente não encontrado.")
		}

		logrus.Error(res.Err())
		return res.Err()
	}

	client := Client{}

	if err := res.Decode(&client); err != nil {
		logrus.Error(err)
		return err
	}

	if err != nil {
		logrus.Error(err)
		return err
	}

	tc := db.Collection("transactions")

	filter = bson.M{"clientid": client.Id}

	opt := options.FindOptions{}
	opt.SetSort(bson.M{"date": -1})
	opt.SetLimit(10)

	cursor, errFind := tc.Find(context.TODO(), filter, &opt)

	if errFind != nil {
		logrus.Error(errFind)
		return errFind
	}

	transactions := []LastTransaction{}
	if err := cursor.All(context.TODO(), &transactions); err != nil {
		logrus.Error(err)
		return err
	}

	statement := Statement{
		Balance: Balance{
			Total: client.Balance,
			Limit: client.Limit,
			Date:  time.Now(),
		},
		LastTransactions: transactions,
	}

	return c.JSON(statement)
}
