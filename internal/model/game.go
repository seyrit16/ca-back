package model

import (
	"ca/internal/model/enum"
	"time"

	"github.com/google/uuid"
)

type Game struct {
	Id uuid.UUID `db:"id"`

	OwnerId *uuid.UUID `db:"owner_id"`

	Status enum.GameStatus `db:"status"`

	WinnerId    uuid.UUID `db:"winner_id"`
	RatingDelta int16     `db:"rating_delta"`

	BoardRadius int16  `db:"board_radius"`
	Seeds       string `db:"seeds"`

	TurnCount           int    `db:"turn_count"`
	CurrentTurnPlayerId string `db:"current_turn_player_id"`
	NumberCharacters    int    `db:"number_characters"`
	NumberCubes         int    `db:"number_cubes"`

	CurrentDisaster   enum.DisasterType `db:"current_disaster"`
	DisasterStartedAt int               `db:"disaster_started_at"`
	DisasterZone      int16             `db:"disaster_zone"`

	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	FinishedAt *time.Time `db:"finished_at"`

	Board   map[string]*GameCell   `db:"-"`
	Players map[string]*GamePlayer `db:"-"`
}
