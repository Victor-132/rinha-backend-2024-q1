package repository

import (
	"context"

	"github.com/Victor-132/rinha-backend2/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ITransactionRepository interface {
	CreateTransaction(t database.Transaction) error
	GetTransactionsByClientId(id int) ([]database.Transaction, error)
}

type TransactionRepository struct {
	db database.IDatabase
}

func NewTransactionRepository(db database.IDatabase) ITransactionRepository {
	return &TransactionRepository{db}
}

func (tr *TransactionRepository) CreateTransaction(t database.Transaction) error {
	tc := tr.db.GetCollection("transactions")

	_, err := tc.InsertOne(context.TODO(), t, nil)
	return err
}

func (tr *TransactionRepository) GetTransactionsByClientId(id int) (t []database.Transaction, err error) {
	tc := tr.db.GetCollection("transactions")

	fil := bson.M{"clientid": id}

	opt := options.FindOptions{}
	opt.SetSort(bson.M{"date": -1})
	opt.SetLimit(10)

	var cur *mongo.Cursor
	cur, err = tc.Find(context.TODO(), fil, &opt)
	if err != nil {
		return
	}

	err = cur.All(context.TODO(), &t)

	return
}
