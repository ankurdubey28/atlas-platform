package main

import (
	"github.com/ankurdubey28/atlas-platform/internal/db"
	"github.com/ankurdubey28/atlas-platform/internal/env"
	"github.com/ankurdubey28/atlas-platform/internal/store"
)

func main() {

	dbCfg := dbConfig{
		addr:         env.GetString("DATABASE_URL", ""),
		maxOpenConns: 5,
		maxIdleConns: 5,
		maxIdleTime:  "30",
	}
	cfg := Config{
		addr:    ":3000",
		db:      dbCfg,
		env:     env.GetString("ENV", "dev"),
		apiURL:  env.GetString("EXTERNAL_URL", "localhost:8080"),
		version: env.GetString("VERSION", "v1"),
	}

	err := db.Connect()
	if err != nil {
		return
	}

	application := app{
		config: cfg,
		store:  store.NewStorage(db.DB),
	}
	mux := application.mount()
	err = application.run(mux)
	if err != nil {
		return
	}
}
