package repository

import (
	"context"
	"errors"

	"github.com/Victor-132/rinha-backend2/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IClientRepository interface {
	GetClientById(id int) (*database.Client, error)
	SetClientBalanceByClientId(id, val int) error
}

type ClientRepository struct {
	db database.IDatabase
}

func NewClientRepository(db database.IDatabase) IClientRepository {
	return &ClientRepository{db}
}

func (cr *ClientRepository) GetClientById(id int) (cl *database.Client, err error) {
	cc := cr.db.GetCollection("clients")

	fil := bson.M{
		"id": id,
	}

	res := cc.FindOne(context.TODO(), fil, nil)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return
		}

		err = res.Err()

		return
	}

	err = res.Decode(&cl)
	if err != nil {
		return
	}

	return
}

func (cr *ClientRepository) SetClientBalanceByClientId(id, val int) (err error) {
	fil := bson.M{
		"id": id,
	}

	upd := bson.M{
		"$inc": bson.M{
			"balance": val,
		},
	}

	cc := cr.db.GetCollection("clients")

	_, err = cc.UpdateOne(context.TODO(), fil, upd, nil)
	return
}
