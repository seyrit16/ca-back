package model

import "github.com/google/uuid"

type GamePlayer struct {
	// PRIMARY KEY (game_id, player_id)
	GameId uuid.UUID `db:"game_id"`
	// PRIMARY KEY (game_id, player_id)
	PlayerId string `db:"player_id"` // unitId внутри игры "player_0" вместо 0 стоит seat (1 или 0)

	UserId uuid.UUID `db:"user_id"`

	Seat int16 `db:"seat"`

	IsWinner bool `db:"is_winner"`

	Units   map[string]*GameUnit `db:"-"`
	Storage []*GameItem          `db:"-"`
}
