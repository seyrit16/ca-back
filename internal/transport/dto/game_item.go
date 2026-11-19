package dto

import (
	"ca/internal/model"
	"ca/internal/model/enum"
)

type GameItemDTO struct {
	Id         int32                 `json:"id"`
	TemplateId int16                 `json:"templateId"`
	Location   enum.ItemLocationType `json:"location"`
	UnitId     string                `json:"unitId,omitempty"`
}

func ToGameItemDTO(item *model.GameItem) *GameItemDTO {
	return &GameItemDTO{
		Id:         item.Id,
		TemplateId: item.TemplateId,
		Location:   item.Location,
		UnitId:     item.UnitId,
	}
}

func GameItemsToDTO(items []*model.GameItem) []*GameItemDTO {
	dtos := make([]*GameItemDTO, len(items))
	for i, item := range items {
		dtos[i] = ToGameItemDTO(item)
	}
	return dtos
}
