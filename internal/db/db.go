package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool crée un pool de connexions PostgreSQL (Supabase)
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(ctx, config)
}
