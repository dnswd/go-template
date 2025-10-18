package infra

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dnswd/arus/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Infra struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*Infra, error) {
	slog.InfoContext(ctx, "Initializing infrastructure...")

	infra := &Infra{}

	if cfg.Database != "" {
		pool, err := NewPostgresPool(ctx, cfg.Database)
		if err != nil {
			return nil, fmt.Errorf("failed to create postgres pool: %w", err)
		}
		infra.pool = pool
		slog.InfoContext(ctx, "Database connection established")
	} else {
		slog.InfoContext(ctx, "Skipping database init")
	}

	return infra, nil
}

func (i *Infra) DB() *pgxpool.Pool {
	return i.pool
}

func (i *Infra) Close() error {
	if i.pool != nil {
		slog.Info("Closing database connection...")
		i.pool.Close()
	}
	return nil
}
