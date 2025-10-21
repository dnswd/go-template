package user

import (
	"context"
	"fmt"

	"github.com/cockroachdb/apd/v3"
	"github.com/google/uuid"
)

type service struct {
	repo Repository // Depends on interface, not implementation
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, email, name, balance string) (*User, error) {
	balanceValue, _, err := apd.NewFromString(balance)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	
	user := &User{
		Email:   email,
		Name:    name,
		Balance: balanceValue,
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

func (s *service) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.FindByID(ctx, id)
}
func (s *service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
