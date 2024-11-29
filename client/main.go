package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	modelA "github.com/real-time-chat/internal/account/model"
	"github.com/real-time-chat/internal/model"
)

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
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatalf("Failed to connect to %v: %v", url, err)
	}
	defer conn.Close()

	var lr = modelA.LoginRequest{
		Email:    "client1@gmail.com",
		Password: "Client1@123",
	}
	err = conn.WriteJSON(&lr)
	if err != nil {
		log.Fatalf("Fail to send json data to server: %v", err)
	}
	var res model.JsonResponse
	err = conn.ReadJSON(&res)
	if err != nil {
		log.Fatalf("Fail to read json response from server: %v", err)
	}
	if res.Error != "" {
		log.Fatalf("Fail to login to server: %v", res.Error)
	}
	return strings.TrimSpace(res.Data["token"])
}

func HandleRegister(url string) {
	registerReq := modelA.AccountRegisterRequest{
		Email:    "client1@gmail.com",
		Username: "client1",
		Password: "Client1@123",
	}
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatalf("Failed to connect to %v: %v", url, err)
	}
	defer conn.Close()

	err = conn.WriteJSON(registerReq)
	if err != nil {
		log.Fatalf("Fail to send register data to server: %v", err)
	}
	var res model.JsonResponse
	err = conn.ReadJSON(&res)
	if err != nil {
		log.Fatalf("Fail to read data from server: %v", res.Error)
	}
	if res.Error != "" {
		log.Fatalf("Fail to register new account: %v", res.Error)
	}
	log.Printf("[%v] %v\n", res.Code, res.Data)
}

func main() {
	// Địa chỉ WebSocket server (thay đổi theo server của bạn)
	serverAddr := "ws://localhost:9909/ws"
	loginAddr := "ws://localhost:9909/login"
	roomAddr := "?room=123"

	// registerAddr := "ws://localhost:9909/register"
	// HandleRegister(registerAddr)
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
			parsedMsg, err := model.UnmarshalMessage(message)
			if err != nil {
				log.Println("Error parse message: ", err)
			}
			content, err := model.ParseContent(parsedMsg.Content)
			if err != nil {
				log.Println("Error parse content: ", err)
			}
			switch content.Type {
			case model.TextType:
				fmt.Printf("%s: %s\n", parsedMsg.User, content.Data)
			case model.MP3FileType:
				fmt.Printf("%s: sent a file\n", parsedMsg.User)
				if err := os.WriteFile("test_receive_file.mp3", content.Data, os.ModePerm); err != nil {
					log.Println(err)
				}
			default:
				fmt.Printf("%s: sent an unknown message", parsedMsg.User)
			}
		}
	}()

	// Đọc tin nhắn từ stdin và gửi đến server
	rootDir, _ := os.Getwd()
	f, err := os.ReadFile(filepath.Join(rootDir, "assets", "test_file_transfer.mp3"))
	if err != nil {
		log.Println("Fail to read file: ", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, model.NewContentBytes(model.MP3FileType, f))
	if err != nil {
		log.Println("Error sending message:", err)
	}
	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	text := scanner.Text()
	// 	if text == "exit" {
	// 		fmt.Println("Exiting chat...")
	// 		break
	// 	}

	// }

	// Đợi goroutine đọc tin nhắn từ server hoàn tất
	<-done
}
