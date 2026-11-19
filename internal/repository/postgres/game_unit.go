package postgres

import (
	"ca/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameUnitRepository interface {
	CreateFromTemplate(ctx context.Context, gameId uuid.UUID, playerId string, templateId int16, unitSeat int16,
		positionQ int, positionR int) (string, error)
	LoadGameUnits(ctx context.Context, gameId uuid.UUID, playerId string) (map[string]*model.GameUnit, error)
}

type gameUnitRepository struct {
	pool *pgxpool.Pool
}

func NewGameUnitRepository(pool *pgxpool.Pool) GameUnitRepository {
	return &gameUnitRepository{pool: pool}
}

func (repo *gameUnitRepository) CreateFromTemplate(ctx context.Context, gameId uuid.UUID, playerId string,
	templateId int16,
	unitSeat int16, positionQ int, positionR int) (string, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	playerSeat := string(playerId[len(playerId)-1])
	unitId := fmt.Sprintf("unit_%s_%d", playerSeat, unitSeat)
	var gameUnitId string
	err = tx.QueryRow(ctx, `
		insert into game_units (game_id, player_id, unit_id, character_template_id, seat, position_q, position_r,
		                        name, hp, max_hp, attack, defense, luck, speed,
		                        critical_chance, critical_damage, accuracy, dodge)
		select $1::uuid, $2, $3, t.id, $5::smallint, $6::smallint, $7::smallint,
		       t.name, t.hp, t.max_hp, t.attack, t.defense, t.luck, t.speed,
		       t.critical_chance, t.critical_damage, t.accuracy, t.dodge
		from character_templates t
		where t.id = $4
		returning unit_id;`,
		gameId, playerId, unitId, templateId, unitSeat, positionQ, positionR).Scan(&gameUnitId)
	if err != nil {
		return "", err
	}

	if err = tx.Commit(ctx); err != nil {

		return "", err
	}

	return gameUnitId, nil
}

func (repo *gameUnitRepository) LoadGameUnits(ctx context.Context, gameId uuid.UUID, playerId string) (map[string]*model.GameUnit, error) {

	units := map[string]*model.GameUnit{}

	query := `select game_id, player_id,unit_id, seat, position_q,position_r,
       character_template_id, name, hp, max_hp,attack,defense, luck,speed,
       critical_chance,critical_damage, accuracy,dodge,equip_head,equip_body,
       equip_boots, equip_weapon,is_alive
		from game_units where game_id = $1 and player_id = $2`

	rows, err := repo.pool.Query(ctx, query, gameId, playerId)
	if err != nil {

		return nil, errors.New("[LoadGameUnits]\n" + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		unit := &model.GameUnit{}

		if err = pgxscan.ScanRow(unit, rows); err != nil {
			return nil, errors.New("[LoadGameUnits]\n" + err.Error())
		}
		units[unit.UnitId] = unit
	}

	return units, nil
}
