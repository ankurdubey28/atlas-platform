package main

import (
	"errors"
	"net/http"

	"github.com/ankurdubey28/atlas-platform/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreatePostPayload struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func (app *app) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := store.User{
		Name:  payload.Name,
		Age:   payload.Age,
		Email: payload.Email,
	}
	err = app.store.Users.Create(r.Context(), &user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	err = app.jsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *app) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.store.Users.GetAll(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	err = app.jsonResponse(w, http.StatusOK, users)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *app) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID

	strVal := chi.URLParam(r, "id")
	if err := id.Scan(strVal); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *app) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID

	strVal := chi.URLParam(r, "id")
	if err := id.Scan(strVal); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err := app.store.Users.Delete(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

}

type UpdateUserPayload struct {
	Name  *string `validate:"omitempty,min=1"`
	Email *string `validate:"omitempty,email"`
	Age   *int    `validate:"omitempty,gte=0"`
}

func (app *app) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var id pgtype.UUID
	val := chi.URLParam(r, "id")

	err := id.Scan(val)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// update or retain user attributes based on what things are supplied in patch

	var payload UpdateUserPayload
	err = readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = validate.Struct(payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Age != nil {
		user.Age = *payload.Age
	}
	if payload.Name != nil {
		user.Name = *payload.Name
	}
	if payload.Email != nil {
		user.Email = *payload.Email
	}

	err = app.store.Users.Update(r.Context(), user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
