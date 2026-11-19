package dto

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UserProfileDataResponse struct {
	Id        pgtype.UUID        `json:"id"`
	Nickname  string             `json:"nickname"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Color     string             `json:"color"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type DraftStateResponse struct {
	YourId      string             `json:"your_id"`
	CurrentTurn string             `json:"current_turn"`
	TurnIndex   int                `json:"turn_index"`
	Picks       map[string][]int16 `json:"picks"`
}
