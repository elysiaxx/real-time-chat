package server

import (
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/real-time-chat/internal/model"
)

type Hub struct {
	rooms      map[string]map[*Client]bool
	Broadcast  chan model.Message
	Register   chan *Subscription
	Unregister chan *Subscription
	mu         sync.Mutex
}

type Subscription struct {
	client *Client
	room   string
}

func NewHub() *Hub {
	var mu sync.Mutex
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan model.Message),
		Register:   make(chan *Subscription),
		Unregister: make(chan *Subscription),
		mu:         mu,
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case sub := <-h.Register:
			h.mu.Lock()
			if h.rooms[string(sub.room[:])] == nil {
				h.rooms[string(sub.room[:])] = make(map[*Client]bool)
			}
			h.rooms[string(sub.room[:])][sub.client] = true
			h.mu.Unlock()

		case sub := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[string(sub.room[:])]; ok {
				delete(clients, sub.client)
				if len(clients) == 0 {
					delete(h.rooms, string(sub.room[:]))
				}
			}
			h.mu.Unlock()

		case msg := <-h.Broadcast:
			h.mu.Lock()
			for client := range h.rooms[msg.Room] {
				select {
				case client.Send <- model.MarshalMessage(&msg):
				default:
					close(client.Send)
					delete(h.rooms[msg.Room], client)
				}
			}
			h.mu.Unlock()
		case <-ticker.C:
			for room := range h.rooms {
				for client := range h.rooms[room] {
					h.CheckToken(room, client)
				}
			}
		}
	}
}

func (h *Hub) CheckToken(room string, c *Client) {
	claims := c.token.Claims.(jwt.MapClaims)
	expiration, ok := claims["exp"].(float64)
	if ok && time.Unix(int64(expiration), 0).Before(time.Now()) {
		// token expired
		c.Conn.WriteMessage(websocket.TextMessage, []byte("Token expired. Disconnecting..."))
		c.Conn.Close()
		delete(h.rooms[room], c)
		return
	}
}
