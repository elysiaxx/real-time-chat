package model

import (
	"encoding/json"
	"time"
)

type Message struct {
	ID       uint //
	Sender   uint // sender id
	Receiver uint // receiver id
	Metadata []byte
	Content  []byte
	Sent_at  *time.Time
	Is_read  bool
}

func Marshal(m *Message) []byte {
	res, _ := json.Marshal(m)
	return res
}

func Unmarshal(data []byte) (*Message, error) {
	var res Message
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
