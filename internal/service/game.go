package service

import (
	"ca/internal/mapgen"
	"ca/internal/repository/postgres"
	"context"
	"log"

	"github.com/google/uuid"
)

type GameService struct {
	gameRepository  postgres.GameRepository
	gameCellService *GameCellService
}

func NewGameService(gameRepository postgres.GameRepository, gameCellService *GameCellService) *GameService {
	return &GameService{gameRepository: gameRepository, gameCellService: gameCellService}
}

func (s *GameService) CreateGameWithPlayer(ctx context.Context, params postgres.CreateGameParams) (uuid.UUID, error) {
	id, err := s.gameRepository.CreateGameWithPlayer(ctx, params)
	if err != nil {
		return uuid.Nil, err
	}

	go func() {
		generator := mapgen.NewMapGenerator(20, 3)
		generator.Generate()
		if err := s.gameCellService.CreateAllGameCellsForGame(context.Background(), id, generator); err != nil {
			log.Println("CreateAllGameCellsForGame error:", err)
		}
	}()

	return id, nil
}
