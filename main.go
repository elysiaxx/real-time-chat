package main

import (
	"context"
	"log"
	"net/http"

	"github.com/real-time-chat/config"
	"github.com/real-time-chat/internal/account"
	"github.com/real-time-chat/internal/database"
	"github.com/real-time-chat/internal/message"
	"github.com/real-time-chat/server"
)

func main() {

	config := config.DefaulConfig()
	ctx := context.Background()
	hub := server.NewHub()
	go hub.Run() // Run the Hub's event loop

	db, err := database.Connect(config.PostgresqlConnectionString)
	if err != nil {
		log.Fatalf("couldn't connect to database: %v", err)
	}

	mongodb, err := database.MongoConnect(config.MongoConnectionString)
	defer func() {
		if err = mongodb.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	mHandler := message.NewHandler(mongodb)
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
	// http.HandleFunc("/transfer", func(w http.ResponseWriter, r *http.Request) {
	// 	server.Transfer(mHandler, aHandler, w, r)
	// })

	log.Println("Socket server started on :", config.SocketServerPort)
	err = http.ListenAndServe(config.SocketServerHost+":"+config.SocketServerPort, nil)
	if err != nil {
		log.Fatal("socket server error: ", err)
	}

	// run http web server
	// fasthttp.ListenAndServe(config.WebServerHost+":"+config.WebServerPort, nil)
}
