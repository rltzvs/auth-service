package entity

import (
	"errors"
	"regexp"
)

var (
	ErrUserAlreadyExists = errors.New("email already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type User struct {
	ID       int    `json:"-"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	if u.Email == "" || u.Password == "" {
		return errors.New("email and password are required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	return nil
}
