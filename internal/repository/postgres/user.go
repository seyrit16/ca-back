package postgres

import (
	"ca/internal/model"
	"ca/internal/transport/dto"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetById(ctx context.Context, id uuid.UUID) (*model.User, error)
	Create(ctx context.Context, data *dto.UserCreateRequest) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (repo *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		select *
		from users
		where username = $1;
	`

	row, err := repo.pool.Query(ctx, query, username)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[model.User])
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) GetById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		select *
		from users
		where id = $1
	`

	var user model.User
	row, err := repo.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	user, err = pgx.CollectOneRow(row, pgx.RowToStructByNameLax[model.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) Create(ctx context.Context, data *dto.UserCreateRequest) (id uuid.UUID, err error) {

	query := `
		insert into users (nickname,username, password,email)
		values ($1, $2, $3, $4)
		returning id;
	`

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	err = repo.pool.QueryRow(ctx, query,
		data.Nickname,
		data.Username,
		hashPassword,
		data.Email).Scan(&id)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (repo *userRepository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	query := `delete from users where id = $1;`

	_, err = repo.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
