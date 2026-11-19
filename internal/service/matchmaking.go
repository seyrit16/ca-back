package service

import (
	"ca/internal/repository/postgres"
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
)

type waitingPlayer struct {
	userId uuid.UUID
	ch     chan uuid.UUID // получит game_id когда найдётся пара
}

type MatchmakingService struct {
	mu          sync.Mutex
	waiting     *waitingPlayer
	channels    map[uuid.UUID]chan uuid.UUID
	gameService *GameService
}

func NewMatchmakingService(gameService *GameService) *MatchmakingService {
	return &MatchmakingService{
		gameService: gameService,
		channels:    make(map[uuid.UUID]chan uuid.UUID),
	}
}

func (s *MatchmakingService) Join(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	s.mu.Lock()
	log.Println("Join called, userId:", userId)
	log.Println("waiting:", s.waiting)

	if s.waiting != nil && s.waiting.userId != userId {
		log.Println("Match found!")
		opponent := s.waiting
		s.waiting = nil
		s.mu.Unlock()

		gameId, err := s.gameService.CreateGameWithPlayer(ctx, postgres.CreateGameParams{
			Player1Id: opponent.userId,
			Player2Id: userId,
		})
		if err != nil {
			return uuid.Nil, err
		}

		opponent.ch <- gameId
		return gameId, nil
	}

	ch := make(chan uuid.UUID, 1)
	player := &waitingPlayer{userId: userId, ch: ch}
	s.waiting = player
	s.channels[userId] = ch
	s.mu.Unlock()

	return uuid.Nil, nil
}

func (s *MatchmakingService) Subscribe(userId uuid.UUID) chan uuid.UUID {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ch, ok := s.channels[userId]; ok {
		return ch
	}

	ch := make(chan uuid.UUID, 1)
	return ch
}

func (s *MatchmakingService) Leave(userId uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.waiting != nil && s.waiting.userId == userId {
		close(s.waiting.ch)
		s.waiting = nil
	}

	delete(s.channels, userId)
}

func (s *MatchmakingService) CleanChannel(userId uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.channels, userId)
}
