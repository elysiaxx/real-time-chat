package server

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (modify for production)
	},
}

func ServeWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	// get token and assign to client
	// do not need to check whether header has authorization and correct token pattern
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	room := r.URL.Query().Get("room")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	client := &Client{
		ID:    r.RemoteAddr,
		Conn:  conn,
		Send:  make(chan []byte),
		token: token,
	}
	sub := &Subscription{
		client: client,
		room:   room,
	}

	hub.Register <- sub

	go client.ReadMessages(hub, sub)
	go client.WriteMessages()
}

func Login(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	w.Header().Set("Content-Type", "application/json")
	type JsonResponse struct {
		Code  int
		Error error
		Data  map[string]string
	}
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var lr LoginRequest
	err = conn.ReadJSON(&lr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if lr.Username == "zboonz" && lr.Password == "abc" {
		token, err := GenerateJWT(lr.Username)
		if err != nil {
			response := JsonResponse{
				Code:  http.StatusInternalServerError,
				Error: err,
				Data:  nil,
			}
			conn.WriteJSON(response)
			return
		}
		response := JsonResponse{Code: http.StatusOK, Error: nil, Data: map[string]string{"token": token}}

		conn.WriteJSON(response)
	} else {
		http.Error(w, "Incorrect username/password", http.StatusBadRequest)
	}
}
