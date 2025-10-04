package user

import "context"

type Repository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	Delete(ctx context.Context, id string) error
}

type Service interface {
	CreateUser(ctx context.Context, email, name string) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
}
