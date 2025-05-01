package websocket

import "sync"

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client] = true
			h.Mutex.Unlock()
		case client := <-h.Unregister:
			h.Mutex.Lock()
			delete(h.Clients, client)
			h.Mutex.Unlock()
		case message := <-h.Broadcast:
			h.Mutex.Lock()
			for c := range h.Clients {
				go c.Write(message)
			}
			h.Mutex.Unlock()
		}
	}
}
