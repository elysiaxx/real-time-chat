package model

import (
	"encoding/json"
	"errors"
)

type ContentType uint

const (
	TextType ContentType = iota
	PDFFileType
	DOCXFileType
	MP3FileType
)

type Message struct {
	Room    string
	User    string
	Content []byte
}

type Content struct {
	Type ContentType
	Data []byte
}

type JsonResponse struct {
	Code  int
	Error string
	Data  map[string]string
}

func MarshalMessage(m *Message) []byte {
	res, _ := json.Marshal(&m)
	return res
}

func UnmarshalMessage(data []byte) (*Message, error) {
	var res Message
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func ParseMessage(data []byte) (*Message, error) {
	if len(data) < 17 {
		return nil, errors.New("too small to parse a message")
	}
	return &Message{
		Room:    string(data[:8]),
		User:    string(data[8:16]),
		Content: data[:],
	}, nil
}

func NewContentBytes(_tp ContentType, data []byte) []byte {
	var res []byte
	res = append(res, byte(_tp))
	return append(res, data...)
}

func ParseContent(content []byte) (*Content, error) {
	if len(content) < 2 {
		return nil, errors.New("too small to Parse to Content")
	}
	return &Content{
		Type: ContentType(content[0]),
		Data: content[1:],
	}, nil
}
