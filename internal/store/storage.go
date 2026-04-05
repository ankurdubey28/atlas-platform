package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Users interface {
		Create(context.Context, *User) error
		GetAll(context.Context) ([]User, error)
		GetByID(context.Context) (*User, error)
		Update(context.Context) error
		Delete(context.Context) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users: &UserStore{db},
	}
}
