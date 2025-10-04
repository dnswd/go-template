package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type service struct {
	repo Repository // Depends on interface, not implementation
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, email, name string) (*User, error) {
	// Business logic here - pure, testable
	user := &User{
		ID:    uuid.NewString(),
		Email: email,
		Name:  name,
	}

	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return user, s.repo.Save(ctx, user)
}

func (s *service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repo.FindByID(ctx, id)
}
func (s *service) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
