package postgres

import (
	"ca/internal/model"
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GamePlayerRepository interface {
	FindPlayerIdByUserId(ctx context.Context, gameId uuid.UUID, userId uuid.UUID) (string, error)
	LoadGamePlayers(ctx context.Context, gameId uuid.UUID) (map[string]*model.GamePlayer, error)
}

type gamePlayerRepository struct {
	pool *pgxpool.Pool
}

func NewGamePlayerRepository(pool *pgxpool.Pool) GamePlayerRepository {
	return &gamePlayerRepository{pool}
}

func (r *gamePlayerRepository) FindPlayerIdByUserId(ctx context.Context, gameId uuid.UUID, userId uuid.UUID) (string,
	error) {

	query := `select player_id from game_players where game_id = $1 and user_id = $2`

	var gamePlayerId string
	err := r.pool.QueryRow(ctx, query, gameId, userId).Scan(&gamePlayerId)
	if err != nil {
		return "", err
	}

	return gamePlayerId, nil
}

func (r *gamePlayerRepository) LoadGamePlayers(ctx context.Context, gameId uuid.UUID) (map[string]*model.GamePlayer, error) {

	players := make(map[string]*model.GamePlayer)

	query := `select game_id, player_id, user_id, seat, is_winner from game_players where game_id = $1`

	rows, err := r.pool.Query(ctx, query, gameId)
	if err != nil {
		return nil, errors.New("[LoadGamePlayers]\n" + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		player := &model.GamePlayer{}
		if err = pgxscan.ScanRow(player, rows); err != nil {
			return nil, errors.New("[LoadGamePlayers]\n" + err.Error())
		}

		players[player.PlayerId] = player
	}

	return players, nil
}
