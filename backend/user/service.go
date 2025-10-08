package user

import (
	"context"
	"fmt"
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
		Email: email,
		Name:  name,
	}

	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	insertedUser, err := s.repo.Save(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to save to db: %w", err)
	}

	return insertedUser, nil
}

func (s *service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repo.FindByID(ctx, id)
}
func (s *service) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
