package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Users interface {
		Create(context.Context, *User) error
		GetAll(context.Context) (*[]User, error)
		GetByID(context.Context, pgtype.UUID) (*User, error)
		Update(context.Context, *User) error
		Delete(context.Context, pgtype.UUID) error
	}
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{
		Users: &UserStore{db},
	}
}
