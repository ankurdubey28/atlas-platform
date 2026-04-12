package main

import (
	"github.com/ankurdubey28/atlas-platform/internal/db"
	"github.com/ankurdubey28/atlas-platform/internal/env"
	"github.com/ankurdubey28/atlas-platform/internal/store"
	"go.uber.org/zap"
)

func main() {

	dbCfg := dbConfig{
		addr:         env.GetString("DATABASE_URL", ""),
		maxOpenConns: 5,
		maxIdleConns: 5,
		maxIdleTime:  "30",
	}
	cfg := Config{
		addr:    env.GetString("PORT", ":3030"),
		db:      dbCfg,
		env:     env.GetString("ENV", "dev"),
		apiURL:  env.GetString("EXTERNAL_URL", "localhost:8080"),
		version: env.GetString("VERSION", "v1"),
	}

	err := db.Connect()
	if err != nil {
		return
	}

	// setup logger
	logger := zap.Must(zap.NewProduction()).Sugar() // sugar logger basically wraps logger to give less verbose api

	application := app{
		config: cfg,
		store:  store.NewStorage(db.DB),
		logger: logger,
	}
	mux := application.mount()
	err = application.run(mux)
	if err != nil {
		return
	}
}
