package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	// Địa chỉ WebSocket server (thay đổi theo server của bạn)
	serverAddr := "ws://localhost:9909/ws"

	// Kết nối tới WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to the WebSocket server!")
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
			fmt.Printf("Server: %s\n", message)
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
