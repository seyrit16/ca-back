package service

import (
	"ca/internal/repository/postgres"
	"context"

	"github.com/google/uuid"
)

type GamePlayerService struct {
	gamePlayerRepository postgres.GamePlayerRepository
}

func NewGamePlayerService(gamePlayerRepository postgres.GamePlayerRepository) *GamePlayerService {
	return &GamePlayerService{gamePlayerRepository}
}

func (s *GamePlayerService) GetPlayerIdByUserId(ctx context.Context, gameId uuid.UUID, userId uuid.UUID) (string,
	error) {
	playerId, err := s.gamePlayerRepository.FindPlayerIdByUserId(ctx, gameId, userId)
	if err != nil {
		return "", err
	}
	return playerId, nil
}
