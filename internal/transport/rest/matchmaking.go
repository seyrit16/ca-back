package rest

import (
	"ca/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MatchmakingHandler struct {
	matchmakingService *service.MatchmakingService
}

func NewMatchmakingHandler(matchmakingService *service.MatchmakingService) *MatchmakingHandler {
	return &MatchmakingHandler{matchmakingService: matchmakingService}
}

// POST /matchmaking/join
func (h *MatchmakingHandler) Join(c *gin.Context) {
	userId := c.MustGet("user_id").(uuid.UUID)

	gameId, err := h.matchmakingService.Join(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if gameId == (uuid.Nil) {
		c.JSON(http.StatusOK, gin.H{"status": "waiting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "found", "game_id": gameId})
}

// DELETE /matchmaking/leave
func (h *MatchmakingHandler) Leave(c *gin.Context) {
	userId := c.MustGet("user_id").(uuid.UUID)

	h.matchmakingService.Leave(userId)
	c.JSON(http.StatusOK, gin.H{"status": "left"})
}

// GET /matchmaking/status  — SSE
func (h *MatchmakingHandler) Status(c *gin.Context) {
	userId := c.MustGet("user_id").(uuid.UUID)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	gameIdCh := h.matchmakingService.Subscribe(userId)

	select {
	case gameId := <-gameIdCh:
		h.matchmakingService.CleanChannel(userId)
		c.SSEvent("match_found", gin.H{"game_id": gameId})
		c.Writer.Flush()
	case <-c.Request.Context().Done():
		h.matchmakingService.Leave(userId)
		return
	}
}
