package user

import (
	"errors"
	"time"
)

// User is the domain entity for a registered user.
type User struct {
	ID        int64     `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Phone     string
	Website   string
	Points    int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// ErrNotFound is returned when a user cannot be found.
var ErrNotFound = errors.New("user not found")

// ErrAlreadyExists is returned when creating a user with a duplicate email.
var ErrAlreadyExists = errors.New("user already exists")

// ErrInsufficientPoints is returned when a user doesn't have enough points for a transfer.
var ErrInsufficientPoints = errors.New("insufficient points")
