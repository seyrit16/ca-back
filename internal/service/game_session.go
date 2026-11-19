package service

import (
	"ca/internal/model"
	"ca/internal/repository/postgres"
	"context"

	"github.com/google/uuid"
)

type GameSessionService struct {
	gameRepository         postgres.GameRepository
	gameCellRepository     postgres.GameCellRepository
	gamePlayerRepository   postgres.GamePlayerRepository
	gameUnitRepository     postgres.GameUnitRepository
	itemTemplateRepository postgres.ItemTemplateRepository
}

func NewGameSessionService(
	gameRepository postgres.GameRepository,
	gameCellRepository postgres.GameCellRepository,
	gamePlayerRepository postgres.GamePlayerRepository,
	gameUnitRepository postgres.GameUnitRepository,
	itemTemplateRepository postgres.ItemTemplateRepository) *GameSessionService {
	return &GameSessionService{
		gameRepository:         gameRepository,
		gameCellRepository:     gameCellRepository,
		gamePlayerRepository:   gamePlayerRepository,
		gameUnitRepository:     gameUnitRepository,
		itemTemplateRepository: itemTemplateRepository,
	}
}

func (s *GameSessionService) GetGameSession(ctx context.Context, gameId uuid.UUID) (*model.Game, error) {
	game, err := s.gameRepository.LoadGame(ctx, gameId)
	if err != nil {
		return nil, err
	}

	game.Board, err = s.gameCellRepository.LoadGameCells(ctx, gameId)
	if err != nil {
		return nil, err
	}

	game.Players, err = s.gamePlayerRepository.LoadGamePlayers(ctx, gameId)
	if err != nil {
		return nil, err
	}

	for _, player := range game.Players {
		player.Units, err = s.gameUnitRepository.LoadGameUnits(ctx, gameId, player.PlayerId)
		if err != nil {
			return nil, err
		}
	}

	return game, nil
}

func (s *GameSessionService) GetItemTemplates(ctx context.Context) (map[int16]*model.ItemTemplate, error) {
	return s.itemTemplateRepository.LoadItemTemplates(ctx)
}
