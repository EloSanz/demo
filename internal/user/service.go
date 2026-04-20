package user

import "context"

type UserService interface {
	GetAll(ctx context.Context, page int, size int) (users []User, total int64, err error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, u User) (User, error)
	Update(ctx context.Context, id int64, u User) (User, error)
	Delete(ctx context.Context, id int64) error
	FetchFromExternal(ctx context.Context) ([]User, error)
	SyncFromExternal(ctx context.Context, externalID int64) (User, error)
	Transfer(ctx context.Context, fromID, toID int64, amount int) error
}
