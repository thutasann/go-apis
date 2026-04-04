package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/thutasann/go-async-api/config"

	_ "github.com/lib/pq"
)

func NewPostgresDb(cfg config.Config) (*sql.DB, error) {
	dsn := cfg.DatabaseUrl()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database connection: %w", err)
	}
	return db, nil
}
