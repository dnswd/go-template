package user

import (
	"context"
	"errors"
	"log"

	"github.com/dnswd/arus/db"
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

func (r *postgresRepo) Save(ctx context.Context, user *User) (*User, error) {
	pgUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email: user.Email,
		Name:  user.Name,
	})

	log.Printf("pgUser %#v", pgUser)

	if err != nil {
		return nil, err
	}

	resultingUser := toUser(pgUser)

	log.Printf("resulting user %#v", resultingUser)
	return resultingUser, resultingUser.Validate()
}

func (r *postgresRepo) FindByID(ctx context.Context, id string) (*User, error) {
	log.Println("ZCZC id", id)
	pgUser, err := r.queries.GetUser(ctx, id)
	log.Println("ZCZC err", err)
	if err != nil {
		return nil, err
	}

	log.Println("ZCZC pgUser", pgUser)

	user := toUser(pgUser)
	log.Println("ZCZC user", user)
	return user, user.Validate()
}

func (r *postgresRepo) Delete(ctx context.Context, id string) error {
	rowsAffected, err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return errors.New("failed to delete id")
	}

	return nil
}
