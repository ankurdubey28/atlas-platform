package api

import (
	"database/sql"

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
		addr:   ":3000",
		db:     dbCfg,
		env:    env.GetString("ENV", "dev"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
	}

	application := app{
		config: cfg,
		store:  store.NewStorage(&sql.DB{}),
	}
	mux := application.mount()
	err := application.run(mux)
	if err != nil {
		return
	}
}
