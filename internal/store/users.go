package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id        pgtype.UUID
	Name      string
	Age       int
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {

}

func (s *UserStore) GetAll(ctx context.Context) ([]User, error) {

}

func (s *UserStore) GetByID(ctx context.Context) (*User, error) {

}
func (s *UserStore) Update(ctx context.Context) error {

}
func (s *UserStore) Delete(ctx context.Context) error {

}
