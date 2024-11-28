package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/real-time-chat/internal/account"
	accountM "github.com/real-time-chat/internal/account/model"
	"github.com/real-time-chat/internal/model"
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

func Register(aH *account.Handler, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	w.Header().Set("Content-Type", "application/json")
	var ARegR accountM.AccountRegisterRequest
	if err = conn.ReadJSON(&ARegR); err != nil {
		log.Printf("Bad request: %s\n", err)
		response := model.JsonResponse{
			Code:  http.StatusInternalServerError,
			Error: fmt.Sprintf("Bad request: %s", err.Error()),
			Data:  make(map[string]string),
		}
		conn.WriteJSON(&response)
		return
	}
	err = aH.Register(&ARegR)
	if err != nil {
		log.Printf("Fail to register: %s\n", err)
		response := model.JsonResponse{
			Code:  http.StatusInternalServerError,
			Error: fmt.Sprintf("Fail to register: %s", err.Error()),
			Data:  make(map[string]string),
		}
		conn.WriteJSON(&response)
		return
	}

	response := model.JsonResponse{
		Code:  http.StatusOK,
		Error: "",
		Data:  map[string]string{"msg": "Create a new account successfully"},
	}
	conn.WriteJSON(&response)
}

func Login(aH *account.Handler, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	w.Header().Set("Content-Type", "application/json")

	var lr accountM.LoginRequest
	err = conn.ReadJSON(&lr)
	if err != nil {
		response := model.JsonResponse{
			Code:  http.StatusInternalServerError,
			Error: fmt.Sprintf("Bad request: %s", err.Error()),
			Data:  nil,
		}

		conn.WriteJSON(response)
		return
	}

	acc, err := aH.Login(&lr)
	if err != nil {

		response := model.JsonResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
			Data:  nil,
		}

		conn.WriteJSON(response)
		return
	}
	token, err := GenerateJWT(acc.Email)
	if err != nil {
		log.Println("Fail to generate jwt token: ", err)
		response := model.JsonResponse{
			Code:  http.StatusInternalServerError,
			Error: fmt.Sprintf("Fail to generate jwt token: %s", err.Error()),
			Data:  nil,
		}

		conn.WriteJSON(response)
		return
	}
	response := model.JsonResponse{Code: http.StatusOK, Error: "", Data: map[string]string{"token": token}}

	conn.WriteJSON(response)
}
