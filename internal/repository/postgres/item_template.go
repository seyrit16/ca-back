package postgres

import (
	"ca/internal/model"
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemTemplateRepository interface {
	LoadItemTemplates(ctx context.Context) (map[int16]*model.ItemTemplate, error)
}

type itemTemplateRepository struct {
	pool *pgxpool.Pool
}

func NewItemTemplateRepository(pool *pgxpool.Pool) *itemTemplateRepository {
	return &itemTemplateRepository{pool: pool}
}

func (r *itemTemplateRepository) LoadItemTemplates(ctx context.Context) (map[int16]*model.ItemTemplate, error) {
	itemTemplates := make(map[int16]*model.ItemTemplate)

	query := `select id,name,description,default_location,created_at,updated_at from item_templates`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, errors.New("[LoadItemTemplates]\n" + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		itemTemplate := &model.ItemTemplate{}

		if err = pgxscan.ScanRow(itemTemplate, rows); err != nil {
			return nil, errors.New("[LoadItemTemplates]\n" + err.Error())
		}
		itemTemplates[itemTemplate.Id] = itemTemplate
	}

	return itemTemplates, nil
}
