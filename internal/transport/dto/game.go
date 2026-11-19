package dto

import (
	"ca/internal/model"
	"ca/internal/model/enum"
	"time"
)

type GameDTO struct {
	Status enum.GameStatus `json:"status"`

	RatingDelta int16 `json:"rating_delta"`

	BoardRadius int16  `json:"board_radius"`
	Seeds       string `json:"seeds"`

	TurnCount           int    `json:"turn_count"`
	CurrentTurnPlayerId string `json:"current_turn_player_id"`
	NumberCharacters    int    `json:"number_characters"`
	NumberCubes         int    `json:"number_cubes"`

	CurrentDisaster   enum.DisasterType `json:"current_disaster"`
	DisasterStartedAt int               `json:"disaster_started_at"`
	DisasterZone      int16             `json:"disaster_zone"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	FinishedAt *time.Time `json:"finished_at"`

	ItemTemplates []*ItemTemplateDTO `json:"item_templates"`

	Board   map[string]*GameCellDTO   `json:"board"`
	Players map[string]*GamePlayerDTO `json:"players"`
}

func ToGameDTO(g *model.Game, t map[int16]*model.ItemTemplate) *GameDTO {
	board := BoardToDTO(g.Board)
	players := GamePlayerToDTO(g.Players)
	item_templates := ItemTemplatesToDTO(t)

	return &GameDTO{
		Status:              g.Status,
		RatingDelta:         g.RatingDelta,
		BoardRadius:         g.BoardRadius,
		Seeds:               g.Seeds,
		TurnCount:           g.TurnCount,
		CurrentTurnPlayerId: g.CurrentTurnPlayerId,
		NumberCharacters:    g.NumberCharacters,
		NumberCubes:         g.NumberCubes,
		CurrentDisaster:     g.CurrentDisaster,
		DisasterStartedAt:   g.DisasterStartedAt,
		DisasterZone:        g.DisasterZone,
		CreatedAt:           g.CreatedAt,
		UpdatedAt:           g.UpdatedAt,
		FinishedAt:          g.FinishedAt,

		ItemTemplates: item_templates,

		Board:   board,
		Players: players,
	}
}
