package user

import (
	"errors"
)

type User struct {
	ID    string
	Email string
	Name  string
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email required")
	}
	return nil
}
