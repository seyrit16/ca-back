package service

import (
	"ca/internal/model"
	"ca/internal/repository/postgres"
	"context"

	"github.com/google/uuid"
)

type GameItemService struct {
	gameItemRepository postgres.GameItemRepository
}

func NewGameItemService(gameItemRepository postgres.GameItemRepository) *GameItemService {
	return &GameItemService{gameItemRepository: gameItemRepository}
}

func (s *GameItemService) CreateFromTemplate(ctx context.Context,
	template *model.ItemTemplate, gameId uuid.UUID,
	playerId string, unitId string) (*model.GameItem, error) {

	item := &model.GameItem{
		GameId:     gameId,
		TemplateId: template.Id,
		PlayerId:   playerId,
		UnitId:     unitId,
		Location:   template.DefaultLocation,
	}

	if err := s.gameItemRepository.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}
