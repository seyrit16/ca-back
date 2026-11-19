package dto

import "ca/internal/model"

type GamePlayerDTO struct {
	PlayerId string `json:"player_id"`

	Seat int16 `json:"seat"`

	Units map[string]*GameUnitDTO `json:"units"`
}

func ToPlayerDTO(player *model.GamePlayer) *GamePlayerDTO {
	units := UnitsToDTO(player.Units)

	return &GamePlayerDTO{
		PlayerId: player.PlayerId,
		Seat:     player.Seat,
		Units:    units,
	}
}

func GamePlayerToDTO(players map[string]*model.GamePlayer) map[string]*GamePlayerDTO {
	result := make(map[string]*GamePlayerDTO, len(players))
	for key, player := range players {
		result[key] = ToPlayerDTO(player)
	}
	return result
}
