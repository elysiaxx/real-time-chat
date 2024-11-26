package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

type JsonResponse struct {
	Code  int
	Error error
	Data  map[string]string
}

func generateJWT() (string, error) {
	jwtSecret := []byte("secret_key")

	// Tạo claims cho token
	claims := jwt.MapClaims{
		"username": "user1",
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // Token hết hạn sau 1 giờ
	}

	// Tạo token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Ký token với bí mật
	return token.SignedString(jwtSecret)
}

func handleLogin(url string) string {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var lr = LoginRequest{
		Username: "zboonz",
		Password: "abc",
	}
	err = conn.WriteJSON(lr)
	if err != nil {
		log.Fatalf("Fail to send json data to server: %v", err)
	}
	var res JsonResponse
	err = conn.ReadJSON(&res)
	if err != nil {
		log.Fatalf("Fail to read json response from server: %v", err)
	}
	if res.Error != nil {
		log.Fatalf("Fail to login to server: %v", res.Error.Error())
	}
	return strings.TrimSpace(res.Data["token"])
}

func main() {
	// Địa chỉ WebSocket server (thay đổi theo server của bạn)
	serverAddr := "ws://localhost:9909/ws"
	loginAddr := "ws://localhost:9909/login"
	roomAddr := "?room=123"
	// Login to server
	token := handleLogin(loginAddr)

	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)

	conn, _, err := websocket.DefaultDialer.Dial(serverAddr+roomAddr, header)
	if err != nil {
		log.Fatalf("Failed to connect to chat server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to the WebSocket chat server!")

	fmt.Println("Type your messages and press Enter to send. Type 'exit' to quit.")

	// Kênh để nhận tin nhắn từ server
	done := make(chan struct{})

	// Goroutine để đọc tin nhắn từ server
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			fmt.Printf("%s\n", message)
		}
	}()

	// Đọc tin nhắn từ stdin và gửi đến server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			fmt.Println("Exiting chat...")
			break
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}

	// Đợi goroutine đọc tin nhắn từ server hoàn tất
	<-done
}
