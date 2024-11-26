package server

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
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
		hub.Broadcast <- Message{room: sub.room, content: msg}
	}
}

func (c *Client) WriteMessages() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, append([]byte(c.GetUsername()+": "), msg...)); err != nil {
			break
		}
	}
}

func (c *Client) GetUsername() string {
	return c.token.Claims.(jwt.MapClaims)["username"].(string)
}
