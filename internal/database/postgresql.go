package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/3n0ugh/allotropes/internal/config"
	"github.com/3n0ugh/allotropes/internal/errors"
	_ "github.com/lib/pq"
)

func OpenConnectionPQ(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.PostgreSQL.DataSource)
	if err != nil {
		return nil, errors.Wrap(err, "postgresql connection")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "ping")
	}

	return db, nil
}
