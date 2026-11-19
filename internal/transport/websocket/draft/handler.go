package draft

import (
	. "ca/internal/transport/websocket"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ca/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type DraftWsHandler struct {
	upgrader          websocket.Upgrader
	gameUnitService   *service.GameUnitService
	gamePlayerService *service.GamePlayerService
	hub               *DraftHub
}

func NewDraftWSHandler(hub *DraftHub, gameUnitService *service.GameUnitService, gamePlayerService *service.GamePlayerService) *DraftWsHandler {
	return &DraftWsHandler{
		hub:               hub,
		gameUnitService:   gameUnitService,
		gamePlayerService: gamePlayerService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *DraftWsHandler) Connect(ctx *gin.Context) {
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

	gameIdStr := ctx.Param("game_id")
	gameId, err := uuid.Parse(gameIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid game_id"})
		return
	}

	conn, err := h.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("draft upgrade error: %v", err)
		return
	}

	client := &Client{
		Conn:   conn,
		UserId: userUUID,
		RoomId: gameIdStr,
		Send:   make(chan []byte, 256),
	}

	h.hub.Register(client)
	go client.WritePump()

	playerId, err := h.gamePlayerService.GetPlayerIdByUserId(ctx.Request.Context(), gameId, userUUID)
	if err != nil {
		conn.Close()
		return
	}

	welcome, _ := json.Marshal(WelcomeMessage{
		Type:     "welcome",
		PlayerId: playerId,
		Seat:     -1,
	})
	client.Send <- welcome

	if err := h.gameUnitService.JoinDraft(gameId, playerId); err != nil {
		conn.Close()
		return
	}

	state, _ := h.gameUnitService.GetDraftState(gameId)
	if state != nil && state.CurrentTurn != "" {
		h.hub.BroadcastDraftState(gameId, state)

		h.startTurnTimer(gameId, 30, func() {
			cp, err := h.gameUnitService.CurrentTurn(ctx.Request.Context(), gameId)
			if err != nil {
				log.Println(err)
			}
			h.autoPick(ctx, gameId, cp)
		})
	}

	defer func() {
		h.hub.Unregister(client)
		conn.Close()
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg IncomingMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "pick":
			h.handlePick(ctx, client, gameId, msg)
		}
	}
}

func (h *DraftWsHandler) startCountdown(gameId uuid.UUID) {
	go func() {
		for i := 10; i >= 0; i-- {
			h.hub.Broadcast(gameId.String(), OutgoingMessage{
				Type:    "countdown",
				Payload: gin.H{"seconds": i},
			})
			time.Sleep(1 * time.Second)
		}

		h.hub.Broadcast(gameId.String(), OutgoingMessage{
			Type:    "draft_complete",
			Payload: nil,
		})
	}()
}

func (h *DraftWsHandler) handlePick(ctx *gin.Context, client *Client, gameId uuid.UUID, msg IncomingMessage) {
	var payload struct {
		TemplateId int16 `json:"template_id"`
	}

	h.cancelTimer(gameId)

	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return
	}

	playerId, err := h.gamePlayerService.GetPlayerIdByUserId(ctx.Request.Context(), gameId, client.UserId)
	if err != nil {
		data, _ := json.Marshal(OutgoingMessage{Type: "error", Payload: gin.H{"message": "player not found"}})
		client.Send <- data
		return
	}

	err = h.gameUnitService.Pick(ctx.Request.Context(), gameId, playerId, payload.TemplateId)
	if err != nil {
		data, _ := json.Marshal(OutgoingMessage{Type: "error", Payload: gin.H{"message": err.Error()}})
		client.Send <- data
		return
	}

	state, err := h.gameUnitService.GetDraftState(gameId)
	if err == nil {
		h.hub.BroadcastDraftState(gameId, state)
	}

	h.startTurnTimer(gameId, 30, func() {
		cp, err := h.gameUnitService.CurrentTurn(ctx.Request.Context(), gameId)
		if err != nil {
			log.Println(err)
		}
		h.autoPick(ctx, gameId, cp)
	})

	if h.gameUnitService.IsDraftComplete(gameId) {
		h.hub.Broadcast(gameId.String(), OutgoingMessage{
			Type:    "draft_full",
			Payload: nil,
		})

		h.cancelTimer(gameId)
		h.startCountdown(gameId)
	}
}

func (h *DraftWsHandler) autoPick(ctx *gin.Context, gameId uuid.UUID, playerId string) {
	err := h.gameUnitService.AutoPick(ctx.Request.Context(), gameId, playerId)
	if err != nil {
		log.Printf("Ошибка в выборе")
		return
	}
	state, _ := h.gameUnitService.GetDraftState(gameId)
	h.hub.BroadcastDraftState(gameId, state)

	h.startTurnTimer(gameId, 30, func() {
		cp, err := h.gameUnitService.CurrentTurn(ctx.Request.Context(), gameId)
		if err != nil {
			log.Println(err)
		}
		h.autoPick(ctx, gameId, cp)
	})

	if h.gameUnitService.IsDraftComplete(gameId) {
		h.hub.Broadcast(gameId.String(), OutgoingMessage{
			Type:    "draft_full",
			Payload: nil,
		})

		h.cancelTimer(gameId)
		h.startCountdown(gameId)
	}
}

func (h *DraftWsHandler) startTurnTimer(gameId uuid.UUID, seconds int, onTimeout func()) {
	h.hub.Lock()
	// отменяем предыдущий таймер если есть
	if cancel, ok := h.hub.timers[gameId.String()]; ok {
		cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	h.hub.timers[gameId.String()] = cancel
	h.hub.Unlock()

	go func() {
		for i := seconds; i >= 0; i-- {
			select {
			case <-ctx.Done():
				return // таймер отменён — игрок сделал пик
			default:
				h.hub.Broadcast(gameId.String(), OutgoingMessage{
					Type:    "countdown",
					Payload: map[string]any{"seconds": i},
				})
				time.Sleep(1 * time.Second)
			}
		}
		onTimeout()
	}()
}

func (h *DraftWsHandler) cancelTimer(gameId uuid.UUID) {
	h.hub.Lock()
	defer h.hub.Unlock()
	if cancel, ok := h.hub.timers[gameId.String()]; ok {
		cancel()
		delete(h.hub.timers, gameId.String())
	}
}
