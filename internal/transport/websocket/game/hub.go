package game

import (
	. "ca/internal/room"
	"ca/internal/service"
	. "ca/internal/transport/websocket"

	"context"
	"log"

	"github.com/google/uuid"
)

type GameHub struct {
	Hub

	rooms      map[string]*Room
	register   chan *Client
	unregister chan *Client

	gameSessionService *service.GameSessionService
}

func NewGameHub(gameSessionService *service.GameSessionService) *GameHub {
	return &GameHub{
		Hub: NewHub(),

		rooms:      make(map[string]*Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),

		gameSessionService: gameSessionService,
	}
}

func (h *GameHub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.Lock()
			room, exists := h.rooms[client.RoomId]
			if !exists {
				failReg := func(err error) {
					log.Println(err)
					close(client.Send)
					h.Unlock()
				}
				var roomId string
				roomId = client.RoomId

				gameId, err := uuid.Parse(roomId)
				if err != nil {
					failReg(err)
				}

				g, err := h.gameSessionService.GetGameSession(ctx, gameId)
				if err != nil {
					failReg(err)
					continue
				}

				it, err := h.gameSessionService.GetItemTemplates(ctx)

				room = NewRoom(client.RoomId, NewState(g, it))
				h.rooms[client.RoomId] = room
				go room.Run(ctx)
			}
			room.AddPlayer(client)
			h.Unlock()

		case client := <-h.unregister:
			h.Lock()
			room, exists := h.rooms[client.RoomId]
			if exists {
				room.RemovePlayer(client)
				if room.IsEmpty() {
					delete(h.rooms, client.RoomId)
				}
			}
			h.Unlock()
		}
	}
}

func (h *GameHub) GetRoom(roomID string) *Room {
	h.RLock()
	defer h.RUnlock()
	return h.rooms[roomID]
}
