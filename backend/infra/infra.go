package infra

import (
	"context"
	"fmt"
	"log"

	"github.com/dnswd/arus/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Infra struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*Infra, error) {
	log.Println("Initializing infrastructure...")

	infra := &Infra{}

	if cfg.Database != "" {
		pool, err := NewPostgresPool(ctx, cfg.Database)
		if err != nil {
			return nil, fmt.Errorf("failed to create postgres pool: %w", err)
		}
		infra.pool = pool
		log.Println("Database connection established")
	}

	return infra, nil
}

func (i *Infra) DB() *pgxpool.Pool {
	return i.pool
}

func (i *Infra) Close() error {
	if i.pool != nil {
		log.Println("Closing database connection...")
		i.pool.Close()
	}
	return nil
}
