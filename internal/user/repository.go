package user

import "context"

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

type UserSearchRepository interface {
	Index(ctx context.Context, u User) error
	Search(ctx context.Context, query string) ([]User, error)
}

type ExternalUserClient interface {
	FetchAll(ctx context.Context) ([]User, error)
	FetchByID(ctx context.Context, id int64) (*User, error)
}
