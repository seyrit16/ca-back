package model

import (
	"ca/internal/model/enum"

	"github.com/google/uuid"
)

type GameItem struct {
	Id int32 `db:"id"`

	GameId uuid.UUID `db:"game_id"`

	TemplateId int16 `db:"template_id"`

	PlayerId string                `db:"player_id"`
	UnitId   string                `db:"unit_id,omitempty"` // nil если на складе
	Location enum.ItemLocationType `db:"location"`
}
