package infra

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pgOnce  sync.Once
	pool    *pgxpool.Pool
	initErr error
)

func NewPostgresPool(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	// This function should only be called once, and won't run second time
	pgOnce.Do(func() {
		log.Println("Initializing PosgreSQL pool...")
		config, err := pgxpool.ParseConfig(connectionString)
		if err != nil {
			initErr = fmt.Errorf("unable to parse PosgreSQL connection string: %v", err)
		}

		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			initErr = fmt.Errorf("unable to connect to PostgreSQL instance: %v", err)
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return pool, nil
}
