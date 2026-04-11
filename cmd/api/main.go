package main

import (
	"context"

	"github.com/ankurdubey28/atlas-platform/internal/env"
	"github.com/ankurdubey28/atlas-platform/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	dbCfg := dbConfig{
		addr:         env.GetString("DATABASE_URL", ""),
		maxOpenConns: 5,
		maxIdleConns: 5,
		maxIdleTime:  "30",
	}
	cfg := Config{
		addr:   ":3000",
		db:     dbCfg,
		env:    env.GetString("ENV", "dev"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
	}

	// create a pool
	pool, err := pgxpool.New(context.Background(), dbCfg.addr)
	if err != nil {
		return
	}

	defer pool.Close()

	application := app{
		config: cfg,
		store:  store.NewStorage(pool),
	}
	mux := application.mount()
	err = application.run(mux)
	if err != nil {
		return
	}
}
