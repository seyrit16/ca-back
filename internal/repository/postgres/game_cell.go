package postgres

import (
	"ca/internal/model"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameCellRepository interface {
	CreateAllGameCellsForGame(ctx context.Context, gameId uuid.UUID, cells map[[2]int16]*model.GameCell) error
	LoadGameCells(ctx context.Context, gameId uuid.UUID) (map[string]*model.GameCell, error)
}

type gameCellRepository struct {
	pool *pgxpool.Pool
}

func NewGameCellRepository(pool *pgxpool.Pool) GameCellRepository {
	return &gameCellRepository{pool: pool}
}

func (r *gameCellRepository) CreateAllGameCellsForGame(ctx context.Context, gameId uuid.UUID,
	cells map[[2]int16]*model.GameCell) error {
	query := "INSERT INTO game_cells (game_id, q, r, generation_group, index_in_group, biome, cell, radius_content) VALUES "

	args := make([]interface{}, 0, len(cells)*8)
	params := make([]string, 0, len(cells))

	idx := 0
	for _, cell := range cells {
		base := idx * 8
		params = append(params, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8))
		args = append(args, gameId, cell.Q, cell.R, cell.GenerationGroup,
			cell.IndexInGroup, cell.Biome, cell.Content, cell.RadiusContent)
		idx++
	}

	query += strings.Join(params, ",")
	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

func (r *gameCellRepository) LoadGameCells(ctx context.Context, gameId uuid.UUID) (map[string]*model.GameCell,
	error) {
	cells := make(map[string]*model.GameCell)

	query := "SELECT game_id, q, r,generation_group,index_in_group,biome,cell,radius_content FROM game_cells WHERE game_id = $1;"

	rows, err := r.pool.Query(ctx, query, gameId)
	if err != nil {
		return nil, errors.New("[LoadGameCells]\n" + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		cell := &model.GameCell{}

		if err = pgxscan.ScanRow(cell, rows); err != nil {
			return nil, errors.New("[LoadGameCells]\n" + err.Error())
		}

		cells[fmt.Sprintf("%d,%d", cell.Q, cell.R)] = cell
	}

	return cells, nil
}
