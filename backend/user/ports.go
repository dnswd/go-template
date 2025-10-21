package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, user *User) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service interface {
	CreateUser(ctx context.Context, email, name, balance string) (*User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
