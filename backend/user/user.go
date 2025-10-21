package user

import (
	"errors"
	"time"

	"github.com/cockroachdb/apd/v3"
	"github.com/google/uuid"
)

type User struct {
	ID        *uuid.UUID
	Email     string
	Name      string
	CreatedAt time.Time
	Balance   *apd.Decimal
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email required")
	}
	if u.Name == "" {
		return errors.New("name required")
	}
	if u.Balance == nil {
		return errors.New("balance required")
	}

	return nil
}
