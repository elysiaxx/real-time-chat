package model

type Message struct {
	Room    string
	Content []byte
}

type JsonResponse struct {
	Code  int
	Error string
	Data  map[string]string
}
