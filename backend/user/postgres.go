package user

import (
	"context"

	"github.com/dnswd/arus/db"
	"github.com/google/uuid"
)

type postgresRepo struct {
	queries *db.Queries
}

func NewPostgresRepository(queries *db.Queries) Repository {
	return &postgresRepo{queries: queries}
}

func toUser(dbUser db.User) *User {
	return &User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt.Time,
	}
}

func (r *postgresRepo) Save(ctx context.Context, user *User) error {
	pgUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:    uuid.NewString(),
		Email: user.Email,
		Name:  user.Name,
	})

	if err != nil {
		return err
	}

	return toUser(pgUser).Validate()
}

func (r *postgresRepo) FindByID(ctx context.Context, id string) (*User, error) {
	pgUser, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	user := toUser(pgUser)
	return user, user.Validate()
}

func (r *postgresRepo) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteUser(ctx, id)
}
