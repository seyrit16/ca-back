package dto

import (
	"ca/internal/model"
	"ca/internal/model/enum"
)

type ItemTemplateDTO struct {
	Id int16 `json:"id"`

	Name        string `json:"itemName"`
	Description string `json:"description"`

	DefaultLocation enum.ItemLocationType `json:"defaultLocation"`
}

func ToItemTemplateDTO(item *model.ItemTemplate) *ItemTemplateDTO {
	return &ItemTemplateDTO{
		Id:              item.Id,
		Name:            item.Name,
		Description:     item.Description,
		DefaultLocation: item.DefaultLocation,
	}
}

func ItemTemplatesToDTO(templates map[int16]*model.ItemTemplate) []*ItemTemplateDTO {
	result := make([]*ItemTemplateDTO, 0, len(templates))
	for _, template := range templates {
		result = append(result, ToItemTemplateDTO(template))
	}
	return result
}
