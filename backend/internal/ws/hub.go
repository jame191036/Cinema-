package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	Hub        *Hub
	Conn       *websocket.Conn
	ShowtimeID string
	Send       chan []byte
}

type Hub struct {
	mu         sync.RWMutex
	rooms      map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *RoomMessage
}

type RoomMessage struct {
	ShowtimeID string
	Data       []byte
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *RoomMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.rooms[client.ShowtimeID] == nil {
				h.rooms[client.ShowtimeID] = make(map[*Client]bool)
			}
			h.rooms[client.ShowtimeID][client] = true
			h.mu.Unlock()
			log.Printf("WS client registered for showtime %s", client.ShowtimeID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.ShowtimeID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.rooms, client.ShowtimeID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.Broadcast:
			h.mu.RLock()
			clients := h.rooms[msg.ShowtimeID]
			for client := range clients {
				select {
				case client.Send <- msg.Data:
				default:
					close(client.Send)
					delete(clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastToRoom(showtimeID string, data []byte) {
	h.Broadcast <- &RoomMessage{ShowtimeID: showtimeID, Data: data}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func ServeWs(hub *Hub, showtimeID string, w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS upgrade error: %v", err)
		return
	}

	client := &Client{
		Hub:        hub,
		Conn:       conn,
		ShowtimeID: showtimeID,
		Send:       make(chan []byte, 256),
	}

	hub.Register <- client
	go client.WritePump()
	go client.ReadPump()
}
