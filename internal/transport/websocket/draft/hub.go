package draft

import (
	. "ca/internal/transport/dto"
	. "ca/internal/transport/websocket"
	"context"

	"github.com/google/uuid"
)

type DraftHub struct {
	Hub
	timers map[string]context.CancelFunc
}

func NewDraftHub() *DraftHub {
	return &DraftHub{
		Hub:    NewHub(),
		timers: make(map[string]context.CancelFunc),
	}
}

func (h *DraftHub) BroadcastDraftState(gameId uuid.UUID, state *DraftStateResponse) {
	h.Broadcast(gameId.String(), OutgoingMessage{
		Type:    "draft_state",
		Payload: state,
	})
}
