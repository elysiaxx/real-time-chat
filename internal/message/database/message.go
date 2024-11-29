package database

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IDbMessage interface {
	Create() error
}

type DbMessage struct {
	db *mongo.Client
}

var _ IDbMessage = (*DbMessage)(nil)

func NewDbMessage(_db *mongo.Client) *DbMessage {
	return &DbMessage{
		db: _db,
	}
}

func (dm *DbMessage) Create() error {
	return nil
}
