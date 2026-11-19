package model

import (
	"ca/internal/model/enum"
	"time"
)

type ItemTemplate struct {
	Id int16 `db:"id"`

	Name        string `db:"name"`
	Description string `db:"description"`

	DefaultLocation enum.ItemLocationType `db:"default_location"`

	Deleted   bool      `db:"deleted"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}
