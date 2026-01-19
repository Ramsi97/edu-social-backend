package websocket

import "sync"

type Client interface {
	ID() string
	Send(message []byte)
	Close() error
}

type Hub struct {
	rooms      map[string]map[string]Client
	register   chan subscription
	unregister chan subscription
	broadcast  chan roomMessage
	mu         sync.RWMutex
}

type subscription struct {
	roomID string
	client Client
}

type roomMessage struct {
	roomID  string
	content []byte
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[string]Client),
		register:   make(chan subscription),
		unregister: make(chan subscription),
		broadcast:  make(chan roomMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.register:
			h.mu.Lock()
			if h.rooms[sub.roomID] == nil {
				h.rooms[sub.roomID] = make(map[string]Client)
			}
			h.rooms[sub.roomID][sub.client.ID()] = sub.client
			h.mu.Unlock()

		case sub := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.rooms[sub.roomID]; ok {
				delete(h.rooms[sub.roomID], sub.client.ID())
				if len(h.rooms[sub.roomID]) == 0 {
					delete(h.rooms, sub.roomID)
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			if clients, ok := h.rooms[msg.roomID]; ok {
				for _, client := range clients {
					client.Send(msg.content)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) JoinRoom(roomID string, client Client) {
	h.register <- subscription{roomID, client}
}

func (h *Hub) BroadcastToRoom(roomID string, message []byte) {
	h.broadcast <- roomMessage{roomID, message}
}