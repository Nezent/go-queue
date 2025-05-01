package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn  *websocket.Conn
	Send  chan []byte
	Mutex sync.Mutex
}

func (c *Client) Write(msg []byte) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Conn.WriteMessage(websocket.TextMessage, msg)
}
