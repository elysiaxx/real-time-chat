package message

import (
	"github.com/real-time-chat/internal/message/database"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handler struct {
	mDB *database.DbMessage
}

func NewHandler(_db *mongo.Client) *Handler {
	return &Handler{mDB: database.NewDbMessage(_db)}
}

func (h *Handler) New() error {
	return nil
}
