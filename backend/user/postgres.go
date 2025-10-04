package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Save(ctx context.Context, user *User) error {
	// Database implementation
	return nil
}

func (r *postgresRepo) FindByID(ctx context.Context, id string) (*User, error) {
	// Database implementation
	return nil, nil
}

func (r *postgresRepo) Delete(ctx context.Context, id string) error {
	// Database implementation
	return nil
}
