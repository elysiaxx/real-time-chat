package database

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

func MongoConnect(url string) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoDB(cl *mongo.Client) *MongoDB {
	return &MongoDB{
		client: cl,
	}
}
