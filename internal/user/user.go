package user

import (
	"errors"
	"time"
)

// User is the domain entity for a registered user.
type User struct {
	ID        int64
	Name      string
	Email     string
	Phone     string
	Website   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ErrNotFound is returned when a user cannot be found.
var ErrNotFound = errors.New("user not found")

// ErrAlreadyExists is returned when creating a user with a duplicate email.
var ErrAlreadyExists = errors.New("user already exists")
