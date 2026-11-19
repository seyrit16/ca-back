package model

import (
	"ca/internal/model/enum"

	"github.com/google/uuid"
)

type GameCell struct {
	GameId uuid.UUID `db:"game_id"`
	Q      int16     `db:"q"`
	R      int16     `db:"r"`

	GenerationGroup *int16 `db:"generation_group"`
	IndexInGroup    *int16 `db:"index_in_group"`

	Biome   enum.BiomeType   `db:"biome"`
	Content enum.ContentType `db:"cell"`

	RadiusContent *int16 `db:"radius_content"`
}
