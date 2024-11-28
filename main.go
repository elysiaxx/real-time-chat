package main

import (
	"log"
	"net/http"

	"github.com/real-time-chat/internal/account"
	"github.com/real-time-chat/internal/database"
	"github.com/real-time-chat/server"
)

const dns = "host=localhost user=postgres password=admin@123 dbname=real-time-chat port=5432 sslmode=disable"

func main() {
	hub := server.NewHub()
	go hub.Run() // Run the Hub's event loop

	db, err := database.Connect(dns)
	if err != nil {
		log.Fatalf("couldn't connect to database: %v", err)
	}
	aHandler := account.NewHandler(db)

	http.HandleFunc("/ws", server.ValidateToken(func(w http.ResponseWriter, r *http.Request) {
		server.ServeWebSocket(hub, w, r)
	}))
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		server.Login(aHandler, w, r)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		server.Register(aHandler, w, r)
	})

	log.Println("Server started on :9909")
	err = http.ListenAndServe(":9909", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
