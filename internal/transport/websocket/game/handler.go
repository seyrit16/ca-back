package game

import (
	websocket2 "ca/internal/transport/websocket"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GameWsHandler struct {
	upgrader websocket.Upgrader
	hub      *GameHub
	logger   *log.Logger
}

func NewGameWSHandler(hub *GameHub) *GameWsHandler {

	return &GameWsHandler{hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,

			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (handler *GameWsHandler) Connect(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userUUID, ok := userId.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	gameID := ctx.Param("game_id")
	if gameID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "game_id required"})
		return
	}

	conn, err := handler.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		handler.logger.Printf("upgrade error: %v", err)
		return
	}

	client := &websocket2.Client{
		Conn:   conn,
		UserId: userUUID,
		RoomId: gameID,
		Send:   make(chan []byte, 256),
	}

	handler.hub.register <- client

	go client.WritePump()

	defer func() {
		handler.hub.unregister <- client
		conn.Close()
	}()

	// читаем сообщения
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				handler.logger.Printf("read error: %v", err)
			}
			break
		}

		var msg websocket2.IncomingMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			handler.logger.Printf("invalid message: %v", err)
			continue
		}
		msg.Raw = data

		room := handler.hub.GetRoom(client.RoomId)
		if room == nil {
			handler.logger.Printf("room not found: %s", client.RoomId)
			break
		}

		room.Actions <- PlayerAction{
			Client:  client,
			Message: msg,
		}
	}
}
