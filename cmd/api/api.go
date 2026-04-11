package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ankurdubey28/atlas-platform/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type app struct {
	config Config
	store  *store.Storage
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

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error running server: %s", err.Error())
	}
	return nil
}
