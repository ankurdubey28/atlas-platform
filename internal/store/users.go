package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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
	db *pgxpool.Pool
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (name, age, email)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	return s.db.QueryRow(ctx, query,
		user.Name,
		user.Age,
		user.Email,
	).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
}

func (s *UserStore) GetAll(ctx context.Context) (*[]User, error) {
	query := `
     SELECT * FROM users
`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var Users []User

	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Name, &u.Age, &u.Email, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		Users = append(Users, u)
	}
	rows.Close()
	return &Users, nil
}

func (s *UserStore) GetByID(ctx context.Context, id pgtype.UUID) (*User, error) {
	query := `
		SELECT id, name, age, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Name,
		&user.Age,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET name = $1,
		    age = $2,
		    email = $3,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`

	cmdTag, err := s.db.Exec(ctx, query,
		user.Name,
		user.Age,
		user.Email,
		user.Id,
	)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found")
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, id pgtype.UUID) error {
	query := `
     DELETE FROM users
     WHERE id=$1
`
	err := s.db.QueryRow(ctx, query, id)
	if err != nil {
		return pgx.ErrNoRows
	}
	return nil
}
