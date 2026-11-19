package dto

import "ca/internal/model"

type GameCellDTO struct {
	Q int16 `json:"q"`
	R int16 `json:"r"`

	GenerationGroup *int16 `json:"generation_group"`
	IndexInGroup    *int16 `json:"index_in_group"`

	Biome   int `json:"biome"`
	Content int `json:"content"`

	RadiusContent *int16 `json:"radius_content"`
}

func ToGameCellDTO(cell *model.GameCell) *GameCellDTO {
	return &GameCellDTO{
		Q:               cell.Q,
		R:               cell.R,
		GenerationGroup: cell.GenerationGroup,
		IndexInGroup:    cell.IndexInGroup,
		Biome:           int(cell.Biome),
		Content:         int(cell.Content),
		RadiusContent:   cell.RadiusContent,
	}
}

func BoardToDTO(board map[string]*model.GameCell) map[string]*GameCellDTO {
	result := make(map[string]*GameCellDTO, len(board))
	for key, cell := range board {
		result[key] = ToGameCellDTO(cell)
	}
	return result
}
