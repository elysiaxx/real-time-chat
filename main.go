package main

import (
	"log"
	"net/http"

	"github.com/real-time-chat/server"
)

func main() {
	hub := server.NewHub()
	go hub.Run() // Run the Hub's event loop

	http.HandleFunc("/ws", server.ValidateToken(func(w http.ResponseWriter, r *http.Request) {
		server.ServeWebSocket(hub, w, r)
	}))

	log.Println("Server started on :9909")
	err := http.ListenAndServe(":9909", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
