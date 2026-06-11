package main

import (
	"context"
	"errors"

	"github.com/ankurdubey28/atlas-platform/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type fakeUserStore struct {
	createFn  func(context.Context, *store.User) error
	getAllFn  func(context.Context) (*[]store.User, error)
	getByIDFn func(context.Context, pgtype.UUID) (*store.User, error)
	updateFn  func(context.Context, *store.User) error
	deleteFn  func(context.Context, pgtype.UUID) error
}

func (f *fakeUserStore) Create(ctx context.Context, user *store.User) error {
	if f.createFn == nil {
		return errors.New("unexpected Create call")
	}
	return f.createFn(ctx, user)
}

func (f *fakeUserStore) GetAll(ctx context.Context) (*[]store.User, error) {
	if f.getAllFn == nil {
		return nil, errors.New("unexpected GetAll call")
	}
	return f.getAllFn(ctx)
}

func (f *fakeUserStore) GetByID(ctx context.Context, id pgtype.UUID) (*store.User, error) {
	if f.getByIDFn == nil {
		return nil, errors.New("unexpected GetByID call")
	}
	return f.getByIDFn(ctx, id)
}

func (f *fakeUserStore) Update(ctx context.Context, user *store.User) error {
	if f.updateFn == nil {
		return errors.New("unexpected Update call")
	}
	return f.updateFn(ctx, user)
}

func (f *fakeUserStore) Delete(ctx context.Context, id pgtype.UUID) error {
	if f.deleteFn == nil {
		return errors.New("unexpected Delete call")
	}
	return f.deleteFn(ctx, id)
}

func newTestApp(users *fakeUserStore) *app {
	return &app{
		config: Config{
			env:     "test",
			version: "test-version",
		},
		store: &store.Storage{
			Users: users,
		},
		logger: zap.NewNop().Sugar(),
	}
}
