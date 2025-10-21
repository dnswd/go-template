package user

import (
	"context"
	"errors"

	"github.com/dnswd/arus/db"
	"github.com/dnswd/arus/util"
	"github.com/google/uuid"
)

type postgresRepo struct {
	queries *db.Queries
}

func NewPostgresRepository(queries *db.Queries) Repository {
	return &postgresRepo{queries: queries}
}

func toUser(dbUser db.User) (*User, error) {
	id, err := util.PgtypeToUUID(dbUser.ID)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        &id,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt.Time,
		Balance:   util.PgtypeToDecimal(dbUser.Balance),
	}, nil
}

func (r *postgresRepo) Save(ctx context.Context, user *User) (*User, error) {
	pgUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:   user.Email,
		Name:    user.Name,
		Balance: util.DecimalToPgtype(user.Balance),
	})

	if err != nil {
		return nil, err
	}

	resultingUser, err := toUser(pgUser)
	if err != nil {
		return nil, err
	}

	return resultingUser, resultingUser.Validate()
}

func (r *postgresRepo) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	pgUser, err := r.queries.GetUser(ctx, util.UUIDToPgtype(id))
	if err != nil {
		return nil, err
	}

	user, err := toUser(pgUser)
	if err != nil {
		return nil, err
	}

	return user, user.Validate()
}

func (r *postgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := r.queries.DeleteUser(ctx, util.UUIDToPgtype(id))
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return errors.New("failed to delete id")
	}

	return nil
}
