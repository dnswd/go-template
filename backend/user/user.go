package user

import (
	"errors"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email required")
	}
	return nil
}
