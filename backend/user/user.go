package user

import (
	"errors"
	"time"
)

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email required")
	}
	return nil
}
