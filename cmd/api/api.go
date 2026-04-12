package main

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/ankurdubey28/atlas-platform/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type app struct {
	config Config
	store  *store.Storage
	logger *zap.SugaredLogger
}

type Config struct {
	addr    string
	db      dbConfig
	env     string
	apiURL  string
	version string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *app) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/health", app.healthCheckHandler)

	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/", app.createUserHandler)
		r.Get("/", app.getAllUsersHandler)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.getUserByIdHandler)
			r.Patch("/", app.updateUserHandler)
			r.Delete("/", app.deleteUserHandler)
		})
	})

	return r
}

func (app *app) run(mux http.Handler) error {
	server := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	ln, err := net.Listen("tcp", app.config.addr)
	if err != nil {
		app.logger.Error("Error Launching Server")
		return err
	}
	app.logger.Infof("Server started on %s", app.config.addr)
	err = server.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return err
	}
	app.logger.Infof("Server shutting down gracefully")
	return nil
}
