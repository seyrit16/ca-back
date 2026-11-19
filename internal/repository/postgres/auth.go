package postgres

import (
	"ca/internal/model"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	GetUserByUsername(ctx context.Context, email string) (*model.User, error)
	SaveRefreshSession(ctx context.Context, session *model.RefreshSession) error
	GetRefreshSession(ctx context.Context, tokenHash string) (*model.RefreshSession, error)
	DeleteRefreshSession(ctx context.Context, tokenHash string) error
	DeleteUserSessions(ctx context.Context, userID int64) error
}

type authRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) AuthRepository {
	return &authRepository{pool: pool}
}

func (repo *authRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	row, err := repo.pool.Query(ctx, `SELECT * FROM users WHERE username = $1`, username)
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
	return &user, nil
}

func (repo *authRepository) SaveRefreshSession(ctx context.Context, session *model.RefreshSession) error {
	_, err := repo.pool.Exec(ctx,
		`INSERT INTO refresh_sessions (user_id, token, expires_at, created_at, user_agent, ip)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (user_id) DO UPDATE
		 SET token = $2, expires_at = $3, created_at = $4, user_agent = $5, ip = $6`,
		session.UserID, session.Token, session.ExpiresAt, time.Now(), session.UserAgent, session.IP,
	)
	return err
}

func (repo *authRepository) GetRefreshSession(ctx context.Context, tokenHash string) (*model.RefreshSession, error) {
	var session model.RefreshSession
	row, err := repo.pool.Query(ctx, `SELECT * FROM refresh_sessions WHERE token = $1`, tokenHash)
	if err != nil {
		return nil, err
	}

	session, err = pgx.CollectOneRow(row, pgx.RowToStructByNameLax[model.RefreshSession])
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (repo *authRepository) DeleteRefreshSession(ctx context.Context, tokenHash string) error {
	_, err := repo.pool.Exec(ctx, `DELETE FROM refresh_sessions WHERE token = $1`, tokenHash)
	return err
}

func (repo *authRepository) DeleteUserSessions(ctx context.Context, userID int64) error {
	_, err := repo.pool.Exec(ctx, `DELETE FROM refresh_sessions WHERE user_id = $1`, userID)
	return err
}
