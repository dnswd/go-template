package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	saveFn func(ctx context.Context, user *User) error
}

func (m *mockRepo) Save(ctx context.Context, user *User) error {
	return m.saveFn(ctx, user)
}

func (m *mockRepo) FindByID(ctx context.Context, id string) (*User, error) {
	return nil, nil
}

func (m *mockRepo) Delete(ctx context.Context, is string) error {
	return nil
}

func TestCreateUser(t *testing.T) {
	// Test business logic WITHOUT database
	mock := &mockRepo{
		saveFn: func(ctx context.Context, user *User) error {
			return nil
		},
	}

	svc := NewService(mock)
	user, err := svc.CreateUser(context.Background(), "test@example.com", "Test")

	assert.NotNil(t, user)
	assert.Nil(t, err)

	// Pure business logic test
}
