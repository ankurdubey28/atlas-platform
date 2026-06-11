package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ankurdubey28/atlas-platform/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestCreateUserHandlerSuccess(t *testing.T) {
	t.Parallel()

	users := &fakeUserStore{
		createFn: func(_ context.Context, user *store.User) error {
			if user.Name != "Ankur" {
				t.Fatalf("expected name to be forwarded, got %q", user.Name)
			}
			if user.Age != 24 {
				t.Fatalf("expected age to be forwarded, got %d", user.Age)
			}
			if user.Email != "ankur@example.com" {
				t.Fatalf("expected email to be forwarded, got %q", user.Email)
			}
			if err := user.Id.Scan("550e8400-e29b-41d4-a716-446655440000"); err != nil {
				return err
			}
			now := time.Unix(1700000000, 0).UTC()
			user.CreatedAt = now
			user.UpdatedAt = now
			return nil
		},
	}
	app := newTestApp(users)

	req := httptest.NewRequest(http.MethodPost, "/v1/users/", strings.NewReader(`{"name":"Ankur","age":24,"email":"ankur@example.com"}`))
	rec := httptest.NewRecorder()

	app.createUserHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body struct {
		Data store.User `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Data.Name != "Ankur" || body.Data.Age != 24 || body.Data.Email != "ankur@example.com" {
		t.Fatalf("unexpected response body: %+v", body.Data)
	}
}

func TestCreateUserHandlerRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	app := newTestApp(&fakeUserStore{
		createFn: func(_ context.Context, _ *store.User) error {
			t.Fatal("Create should not be called when payload is invalid")
			return nil
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/users/", strings.NewReader(`{"name":"Ankur","age":24,"email":"ankur@example.com","role":"admin"}`))
	rec := httptest.NewRecorder()

	app.createUserHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetAllUsersHandlerSuccess(t *testing.T) {
	t.Parallel()

	id := pgtype.UUID{}
	if err := id.Scan("550e8400-e29b-41d4-a716-446655440000"); err != nil {
		t.Fatalf("failed to build UUID: %v", err)
	}

	users := []store.User{
		{Name: "Ankur", Age: 24, Email: "ankur@example.com", Id: id},
	}
	app := newTestApp(&fakeUserStore{
		getAllFn: func(_ context.Context) (*[]store.User, error) {
			return &users, nil
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/v1/users/", nil)
	rec := httptest.NewRecorder()

	app.getAllUsersHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestGetUserByIDHandlerRejectsInvalidUUID(t *testing.T) {
	t.Parallel()

	app := newTestApp(&fakeUserStore{})
	req := httptest.NewRequest(http.MethodGet, "/v1/users/not-a-uuid/", nil)
	rec := httptest.NewRecorder()
	req = withURLParam(req, "id", "not-a-uuid")

	app.getUserByIdHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestDeleteUserHandlerReturnsNotFound(t *testing.T) {
	t.Parallel()

	app := newTestApp(&fakeUserStore{
		deleteFn: func(_ context.Context, _ pgtype.UUID) error {
			return pgx.ErrNoRows
		},
	})
	req := httptest.NewRequest(http.MethodDelete, "/v1/users/550e8400-e29b-41d4-a716-446655440000/", nil)
	req = withURLParam(req, "id", "550e8400-e29b-41d4-a716-446655440000")
	rec := httptest.NewRecorder()

	app.deleteUserHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUpdateUserHandlerPartiallyUpdatesExistingUser(t *testing.T) {
	t.Parallel()

	id := pgtype.UUID{}
	if err := id.Scan("550e8400-e29b-41d4-a716-446655440000"); err != nil {
		t.Fatalf("failed to build UUID: %v", err)
	}

	existing := &store.User{
		Id:    id,
		Name:  "Ankur",
		Age:   24,
		Email: "ankur@example.com",
	}

	users := &fakeUserStore{
		getByIDFn: func(_ context.Context, gotID pgtype.UUID) (*store.User, error) {
			if gotID != id {
				t.Fatal("expected requested UUID to match existing user")
			}
			clone := *existing
			return &clone, nil
		},
		updateFn: func(_ context.Context, user *store.User) error {
			if user.Name != "Updated" {
				t.Fatalf("expected updated name, got %q", user.Name)
			}
			if user.Age != 24 {
				t.Fatalf("expected untouched age, got %d", user.Age)
			}
			if user.Email != "ankur@example.com" {
				t.Fatalf("expected untouched email, got %q", user.Email)
			}
			return nil
		},
	}
	app := newTestApp(users)

	req := httptest.NewRequest(http.MethodPatch, "/v1/users/550e8400-e29b-41d4-a716-446655440000/", strings.NewReader(`{"name":"Updated"}`))
	req = withURLParam(req, "id", "550e8400-e29b-41d4-a716-446655440000")
	rec := httptest.NewRecorder()

	app.updateUserHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestUpdateUserHandlerRejectsInvalidPayload(t *testing.T) {
	t.Parallel()

	id := pgtype.UUID{}
	if err := id.Scan("550e8400-e29b-41d4-a716-446655440000"); err != nil {
		t.Fatalf("failed to build UUID: %v", err)
	}

	app := newTestApp(&fakeUserStore{
		getByIDFn: func(_ context.Context, _ pgtype.UUID) (*store.User, error) {
			return &store.User{Id: id, Name: "Ankur", Age: 24, Email: "ankur@example.com"}, nil
		},
		updateFn: func(_ context.Context, _ *store.User) error {
			t.Fatal("Update should not be called for invalid payload")
			return nil
		},
	})
	req := httptest.NewRequest(http.MethodPatch, "/v1/users/550e8400-e29b-41d4-a716-446655440000/", strings.NewReader(`{"age":-1}`))
	req = withURLParam(req, "id", "550e8400-e29b-41d4-a716-446655440000")
	rec := httptest.NewRecorder()

	app.updateUserHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateUserHandlerReturnsInternalServerErrorOnInvalidUUID(t *testing.T) {
	t.Parallel()

	app := newTestApp(&fakeUserStore{})
	req := httptest.NewRequest(http.MethodPatch, "/v1/users/not-a-uuid/", strings.NewReader(`{"name":"Updated"}`))
	req = withURLParam(req, "id", "not-a-uuid")
	rec := httptest.NewRecorder()

	app.updateUserHandler(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func withURLParam(req *http.Request, key, value string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}
