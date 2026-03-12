package ws

import (
	"encoding/json"
	"log"
	"sync"

	fiberws "github.com/gofiber/websocket/v2"
)

// Event is broadcast to all connected war-room clients.
type Event struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type client struct {
	conn *fiberws.Conn
	send chan []byte
}

// Hub manages all connected WebSocket clients.
type Hub struct {
	mu      sync.RWMutex
	clients map[*client]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*client]struct{})}
}

// Broadcast sends an event to every connected client.
func (h *Hub) Broadcast(e Event) {
	data, err := json.Marshal(e)
	if err != nil {
		log.Printf("ws hub: marshal: %v", err)
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		select {
		case c.send <- data:
		default:
			// slow client — drop
		}
	}
}

func (h *Hub) register(c *client)   { h.mu.Lock(); h.clients[c] = struct{}{}; h.mu.Unlock() }
func (h *Hub) unregister(c *client) { h.mu.Lock(); delete(h.clients, c); h.mu.Unlock() }

// Handler returns the Fiber WebSocket handler for /ws/warroom.
func (h *Hub) Handler() func(*fiberws.Conn) {
	return func(conn *fiberws.Conn) {
		c := &client{conn: conn, send: make(chan []byte, 64)}
		h.register(c)
		defer func() {
			h.unregister(c)
			conn.Close()
		}()

		// Keep-alive read (blocks until client disconnects)
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					close(c.send)
					return
				}
			}
		}()

		for msg := range c.send {
			if err := conn.WriteMessage(fiberws.TextMessage, msg); err != nil {
				return
			}
		}
	}
}
