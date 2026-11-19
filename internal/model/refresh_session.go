package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshSession struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"` // хэш refresh токена
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	UserAgent string    `db:"user_agent"`
	IP        string    `db:"ip"`
}
