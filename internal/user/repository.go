package user

import "context"

// UserRepository defines the persistence contract for the User domain.
type UserRepository interface {
	FindAll(ctx context.Context, page int, size int) (users []User, total int64, err error)
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, u User) (User, error)
	Update(ctx context.Context, u User) (User, error)
	Delete(ctx context.Context, id int64) error
	TransferPoints(ctx context.Context, fromID, toID int64, amount int) error
}

// UserSearchRepository defines the contract for fuzzy searching users.
type UserSearchRepository interface {
	Index(ctx context.Context, u User) error
	Search(ctx context.Context, query string) ([]User, error)
}

// ExternalUserClient is the port for fetching users from the JSONPlaceholder external API.
type ExternalUserClient interface {
	FetchAll(ctx context.Context) ([]User, error)
	// FetchByID returns nil when the user does not exist in the external system.
	FetchByID(ctx context.Context, id int64) (*User, error)
}
