package websocket

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[string][]*Client
}

func NewHub() Hub {
	return Hub{clients: make(map[string][]*Client)}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.RoomId] = append(h.clients[client.RoomId], client)
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	clients := h.clients[client.RoomId]
	for i, c := range clients {
		if c == client {
			h.clients[client.RoomId] = append(clients[:i], clients[i+1:]...)
			close(client.Send)
			break
		}
	}
	if len(h.clients[client.RoomId]) == 0 {
		delete(h.clients, client.RoomId)
	}
}

func (h *Hub) Broadcast(roomId string, msg OutgoingMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()
	data, _ := json.Marshal(msg)
	for _, client := range h.clients[roomId] {
		client.Send <- data
	}
}

func (h *Hub) GetClients(roomId string) []*Client {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.clients[roomId]
}

// для mu
func (h *Hub) Lock()    { h.mu.Lock() }
func (h *Hub) Unlock()  { h.mu.Unlock() }
func (h *Hub) RLock()   { h.mu.RLock() }
func (h *Hub) RUnlock() { h.mu.RUnlock() }
