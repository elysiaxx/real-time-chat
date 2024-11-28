package server

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/real-time-chat/internal/model"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte

	token *jwt.Token
}

func (c *Client) ReadMessages(hub *Hub, sub *Subscription) {
	defer func() {
		hub.Unregister <- sub
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		hub.Broadcast <- model.Message{Room: sub.room, Content: msg}
	}
}

func (c *Client) WriteMessages() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, append([]byte(c.GetUser()+": "), msg...)); err != nil {
			break
		}
	}
}

func (c *Client) GetUser() string {
	return c.token.Claims.(jwt.MapClaims)["email"].(string)
}
