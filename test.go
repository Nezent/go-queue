// // Directory: go-queue/

// // ============================
// // cmd/routes/routes.go
package routes

// import (
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"go-queue/internal/handler"
// 	"go-queue/internal/websocket"
// )

// func RegisterRoutes(r chi.Router, hub *websocket.Hub) {
// 	r.Get("/ws/jobs", func(w http.ResponseWriter, r *http.Request) {
// 		handler.HandleWebSocket(hub, w, r)
// 	})
// }

// // ============================
// // cmd/worker/main.go
// package main

// import (
// 	"log"

// 	"go-queue/internal/bootstrap"
// 	"go-queue/internal/db"
// )

// func main() {
// 	hub := bootstrap.SetupWebSocketHub()
// 	go hub.Run()
// 	db.StartPgListener("job_notifications", hub)
// 	select {} // Block forever
// }

// // ============================
// // cmd/server/main.go
// package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"go-queue/cmd/routes"
// 	"go-queue/internal/bootstrap"
// )

// func main() {
// 	r := chi.NewRouter()
// 	hub := bootstrap.SetupWebSocketHub()
// 	go hub.Run()

// 	routes.RegisterRoutes(r, hub)

// 	log.Println("Server started on :8080")
// 	http.ListenAndServe(":8080", r)
// }

// // ============================
// // internal/bootstrap/websocket.go
// package bootstrap

// import "go-queue/internal/websocket"

// func SetupWebSocketHub() *websocket.Hub {
// 	hub := websocket.NewHub()
// 	return hub
// }

// // ============================
// // internal/db/listener.go
// package db

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"go-queue/internal/websocket"
// )

// func StartPgListener(channel string, hub *websocket.Hub) {
// 	ctx := context.Background()
// 	db, err := pgxpool.New(ctx, "postgres://youruser:yourpass@localhost:5432/yourdb")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	conn, err := db.Acquire(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Release()

// 	_, err = conn.Exec(ctx, "LISTEN " + channel)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("Listening for notifications on:", channel)

// 	for {
// 		n, err := conn.Conn().WaitForNotification(ctx)
// 		if err != nil {
// 			log.Println("listen error:", err)
// 			time.Sleep(1 * time.Second)
// 			continue
// 		}
// 		hub.Broadcast <- []byte(n.Payload)
// 	}
// }

// // ============================
// // internal/handler/websocket.go
// package handler

// import (
// 	"net/http"

// 	"github.com/gorilla/websocket"
// 	"go-queue/internal/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// func HandleWebSocket(hub *websocket.Hub, w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		return
// 	}
// 	client := &websocket.Client{
// 		Conn: conn,
// 		Send: make(chan []byte, 256),
// 	}

// 	hub.Register <- client

// 	defer func() {
// 		hub.Unregister <- client
// 		conn.Close()
// 	}()

// 	for {
// 		_, _, err := conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 	}
// }

// // ============================
// // internal/websocket/client.go
// package websocket

// import (
// 	"sync"
// 	"github.com/gorilla/websocket"
// )

// type Client struct {
// 	Conn  *websocket.Conn
// 	Send  chan []byte
// 	Mutex sync.Mutex
// }

// func (c *Client) Write(msg []byte) {
// 	c.Mutex.Lock()
// 	defer c.Mutex.Unlock()
// 	c.Conn.WriteMessage(websocket.TextMessage, msg)
// }

// // ============================
// // internal/websocket/hub.go
// package websocket

// import "sync"

// type Hub struct {
// 	Clients    map[*Client]bool
// 	Broadcast  chan []byte
// 	Register   chan *Client
// 	Unregister chan *Client
// 	Mutex      sync.Mutex
// }

// func NewHub() *Hub {
// 	return &Hub{
// 		Clients:    make(map[*Client]bool),
// 		Broadcast:  make(chan []byte),
// 		Register:   make(chan *Client),
// 		Unregister: make(chan *Client),
// 	}
// }

// func (h *Hub) Run() {
// 	for {
// 		select {
// 		case client := <-h.Register:
// 			h.Mutex.Lock()
// 			h.Clients[client] = true
// 			h.Mutex.Unlock()
// 		case client := <-h.Unregister:
// 			h.Mutex.Lock()
// 			delete(h.Clients, client)
// 			h.Mutex.Unlock()
// 		case message := <-h.Broadcast:
// 			h.Mutex.Lock()
// 			for c := range h.Clients {
// 				go c.Write(message)
// 			}
// 			h.Mutex.Unlock()
// 		}
// 	}
// }
