package user

import (
	"errors"
	"time"
)

type User struct {
	ID        int64  `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Phone     string
	Website   string
	Points    int       `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

var ErrNotFound = errors.New("user not found")

var ErrAlreadyExists = errors.New("user already exists")

var ErrInsufficientPoints = errors.New("insufficient points")
