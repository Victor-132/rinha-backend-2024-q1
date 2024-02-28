package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDatabase interface {
	Connect() error
	Disconnect()
	InitDB() error
	GetCollection(string) *mongo.Collection
}

type Database struct {
	db *mongo.Database
}

func NewDatabase() IDatabase {
	return &Database{}
}

func (d *Database) Connect() error {
	opt := options.Client().
		ApplyURI(os.Getenv("DATABASE_URL"))

	cl, err := mongo.Connect(context.TODO(), opt)
	if err != nil {
		return err
	}

	err = cl.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	d.db = cl.Database("rinha")

	return nil
}

func (d *Database) Disconnect() {
	d.db.Client().Disconnect(context.TODO())
}

func (d *Database) InitDB() (err error) {
	cc := d.db.Collection("clients")
	tc := d.db.Collection("transactions")

	_, err = cc.DeleteMany(context.TODO(), bson.M{}, nil)
	if err != nil {
		return
	}

	_, err = tc.DeleteMany(context.TODO(), bson.M{}, nil)
	if err != nil {
		return
	}

	cls := []Client{
		{Id: 1, Limit: 100000, Balance: 0},
		{Id: 2, Limit: 80000, Balance: 0},
		{Id: 3, Limit: 1000000, Balance: 0},
		{Id: 4, Limit: 10000000, Balance: 0},
		{Id: 5, Limit: 500000, Balance: 0},
	}

	fil := bson.M{}
	upd := bson.M{}

	for _, c := range cls {
		fil["id"] = c.Id
		upd["$set"] = c

		opt := options.UpdateOptions{}
		opt.SetUpsert(true)

		_, err = cc.UpdateOne(context.TODO(), fil, upd, &opt)
		if err != nil {
			return
		}
	}

	_, err = cc.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{Keys: bson.D{{Key: "id", Value: 1}}},
	)
	if err != nil {
		return
	}

	_, err = tc.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{Keys: bson.D{
			{Key: "clientid", Value: 1},
			{Key: "date", Value: 1},
		}},
	)
	if err != nil {
		return
	}

	return
}

func (d *Database) GetCollection(col string) *mongo.Collection {
	return d.db.Collection(col)
}
