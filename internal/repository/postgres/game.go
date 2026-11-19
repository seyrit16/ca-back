package postgres

import (
	"ca/internal/model"
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameRepository interface {
	CreateGameWithPlayer(ctx context.Context, params CreateGameParams) (uuid.UUID, error)
	LoadGame(ctx context.Context, gameId uuid.UUID) (*model.Game, error)
}

type gameRepository struct {
	pool *pgxpool.Pool
}

func NewGameRepository(pool *pgxpool.Pool) GameRepository {
	return &gameRepository{pool: pool}
}

type CreateGameParams struct {
	OwnerId   *uuid.UUID
	Player1Id uuid.UUID
	Player2Id uuid.UUID
}

func (repo *gameRepository) CreateGameWithPlayer(ctx context.Context, params CreateGameParams) (uuid.UUID, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	var gameId uuid.UUID
	err = tx.QueryRow(ctx, `
        INSERT INTO games (owner_id)
        VALUES ($1)
        RETURNING id
    `, params.OwnerId).Scan(&gameId)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO game_players (game_id, player_id, user_id, seat)
        VALUES ($1, 'player_0', $2, 0), ($1, 'player_1', $3, 1)
    `, gameId, params.Player1Id, params.Player2Id)
	if err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return gameId, nil
}

func (repo *gameRepository) LoadGame(ctx context.Context, gameId uuid.UUID) (*model.Game, error) {
	game := &model.Game{}

	query := `
		SELECT id, owner_id, status, winner_id, rating_delta,
       	board_radius, seeds, turn_count,current_turn_player_id,
       	number_characters,number_cubes,current_disaster, disaster_started_at,
       	disaster_zone, created_at, updated_at, finished_at 
		FROM games  WHERE id = $1`

	if err := pgxscan.Get(ctx, repo.pool, game, query, gameId); err != nil {
		return nil, errors.New("[LoadGame]\n" + err.Error())
	}

	return game, nil
}
