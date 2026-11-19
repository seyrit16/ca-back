package postgres

import (
	"ca/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GameItemRepository interface {
	Create(ctx context.Context, item *model.GameItem) error
}

type gameItemRepository struct {
	pool *pgxpool.Pool
}

func NewGameItemRepository(pool *pgxpool.Pool) *gameItemRepository {
	return &gameItemRepository{pool: pool}
}

func (r *gameItemRepository) Create(ctx context.Context, item *model.GameItem) error {
	query := `
        INSERT INTO game_items (id,game_id, template_id, player_id, unit_id, location)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	return r.pool.QueryRow(ctx, query,
		item.Id,
		item.GameId,
		item.TemplateId,
		item.PlayerId,
		item.UnitId,
		item.Location,
	).Scan(&item.Id)
}
