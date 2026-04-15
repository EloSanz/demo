package storage

import (
	"context"
	"errors"
)

// StorageService defines the contract for file storage operations.
type StorageService interface {
	// Upload stores content under the given file name and returns a confirmation message.
	Upload(ctx context.Context, name string, content []byte) (string, error)

	// Download retrieves the raw bytes of a file by name.
	// Returns ErrNotFound when no file with the given name exists.
	Download(ctx context.Context, name string) ([]byte, error)

	// Delete removes a file from the storage backend.
	Delete(ctx context.Context, name string) error

	// List returns all file names in the storage backend.
	List(ctx context.Context) ([]string, error)
}

// ErrNotFound is returned when a requested file does not exist.
var ErrNotFound = errors.New("file not found")

// ErrStorage wraps unexpected storage backend errors.
var ErrStorage = errors.New("storage error")
