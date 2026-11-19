package service

import (
	"ca/internal/repository/postgres"
	. "ca/internal/transport/dto"
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// snake draft: 0-1-1-0-0-1-1-0-0-1
var snakeOrder = []int{0, 1, 1, 0, 0, 1, 1, 0, 0, 1}

type DraftState struct {
	GameId uuid.UUID

	PlayerOrder []string

	CurrentTurn string             // player_id чья очередь
	TurnIndex   int                // 0..9 (5 пиков * 2 игрока)
	Picks       map[string][]int16 // player_id -> []templateId
}

type GameUnitService struct {
	gameUnitRepository postgres.GameUnitRepository
	drafts             map[uuid.UUID]*DraftState // в памяти
	mu                 sync.Mutex
}

func NewGameUnitService(gameUnitRepository postgres.GameUnitRepository) *GameUnitService {
	return &GameUnitService{
		gameUnitRepository: gameUnitRepository,
		drafts:             make(map[uuid.UUID]*DraftState),
	}
}

func (s *GameUnitService) JoinDraft(gameId uuid.UUID, playerId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	draft, ok := s.drafts[gameId]
	if !ok {
		draft = &DraftState{
			GameId: gameId,
			Picks:  make(map[string][]int16),
		}
		s.drafts[gameId] = draft
	}

	if _, exists := draft.Picks[playerId]; exists {
		return errors.New("already joined")
	}

	draft.Picks[playerId] = []int16{}

	if len(draft.Picks) == 2 {
		players := make([]string, 0, 2)
		for p := range draft.Picks {
			players = append(players, p)
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(players), func(i, j int) {
			players[i], players[j] = players[j], players[i]
		})
		draft.PlayerOrder = players
		draft.CurrentTurn = players[snakeOrder[0]]
	}

	return nil
}

func (s *GameUnitService) Pick(ctx context.Context, gameId uuid.UUID, playerId string, templateId int16) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	draft, ok := s.drafts[gameId]
	if !ok {
		return errors.New("draft not found")
	}
	if len(draft.Picks) < 2 {
		return errors.New("waiting for opponent")
	}
	if draft.CurrentTurn != playerId {
		return errors.New("not your turn")
	}
	if len(draft.Picks[playerId]) >= 5 {
		return errors.New("already picked 5 units")
	}

	draft.Picks[playerId] = append(draft.Picks[playerId], templateId)
	seat := int16(len(draft.Picks[playerId]) - 1)
	startPositions := map[string][5][2]int{
		"player_0": {{-18, 18}, {-19, 19}, {-19, 20}, {-20, 19}, {-20, 20}},
		"player_1": {{18, -18}, {19, -19}, {19, -20}, {20, -19}, {20, -20}},
	}
	_, err := s.gameUnitRepository.CreateFromTemplate(ctx, gameId, playerId, templateId, seat,
		startPositions[playerId][seat][0], startPositions[playerId][seat][1])
	if err != nil {
		return errors.New("Error on pick unit")
	}

	draft.TurnIndex++

	if draft.TurnIndex < 10 {
		draft.CurrentTurn = draft.PlayerOrder[snakeOrder[draft.TurnIndex]]
	}

	return nil
}

func (s *GameUnitService) AutoPick(ctx context.Context, gameId uuid.UUID, playerId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	draft, ok := s.drafts[gameId]
	if !ok {
		return errors.New("draft not found")
	}
	if len(draft.Picks) < 2 {
		return errors.New("waiting for opponent")
	}
	if draft.CurrentTurn != playerId {
		return errors.New("not your turn")
	}
	if len(draft.Picks[playerId]) >= 5 {
		return errors.New("already picked 5 units")
	}

	seat := int16(len(draft.Picks[playerId]) - 1)
	startPositions := map[string][5][2]int{
		"player_0": {{-18, 18}, {-19, 19}, {-19, 20}, {-20, 19}, {-20, 20}},
		"player_1": {{18, -18}, {19, -19}, {19, -20}, {20, -19}, {20, -20}},
	}
	_, err := s.gameUnitRepository.CreateFromTemplate(ctx, gameId, playerId, 0, seat,
		startPositions[playerId][seat][0], startPositions[playerId][seat][1])
	if err != nil {
		return errors.New("Error on pick unit")
	}

	draft.TurnIndex++

	if draft.TurnIndex < 10 {
		draft.CurrentTurn = draft.PlayerOrder[snakeOrder[draft.TurnIndex]]
	}

	return nil
}

func (s *GameUnitService) CurrentTurn(ctx context.Context, gameId uuid.UUID) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.drafts[gameId].CurrentTurn, nil
}

func (s *GameUnitService) IsDraftComplete(gameId uuid.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	draft, ok := s.drafts[gameId]
	if !ok {
		return false
	}

	if len(draft.Picks) < 2 {
		return false
	}

	for _, picks := range draft.Picks {
		if len(picks) < 5 {
			return false
		}
	}

	return true
}

func (s *GameUnitService) GetDraftState(gameId uuid.UUID) (*DraftStateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	draft, ok := s.drafts[gameId]
	if !ok {
		return nil, errors.New("draft not found")
	}

	return &DraftStateResponse{
		CurrentTurn: draft.CurrentTurn,
		TurnIndex:   draft.TurnIndex,
		Picks:       draft.Picks,
	}, nil
}
