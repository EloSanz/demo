package user

import "context"

// UserService defines the use-case contract for the User bounded context.
type UserService interface {
	GetAll(ctx context.Context, page int, size int) (users []User, total int64, err error)
	// GetByID returns ErrNotFound when the user does not exist.
	GetByID(ctx context.Context, id int64) (*User, error)
	// Create returns ErrAlreadyExists when the email is already registered.
	Create(ctx context.Context, u User) (User, error)
	// Update returns ErrNotFound when the user does not exist.
	Update(ctx context.Context, id int64, u User) (User, error)
	// Delete returns ErrNotFound when the user does not exist.
	Delete(ctx context.Context, id int64) error
	FetchFromExternal(ctx context.Context) ([]User, error)
	// SyncFromExternal fetches a user from JSONPlaceholder and upserts them locally.
	// Returns ErrNotFound when the external user does not exist.
	SyncFromExternal(ctx context.Context, externalID int64) (User, error)
	// Transfer moves points between two users atomically.
	Transfer(ctx context.Context, fromID, toID int64, amount int) error
}
