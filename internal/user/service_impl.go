package user

import (
	"context"
	"fmt"
)

// userService implements UserService.
type userService struct {
	repo           UserRepository
	externalClient ExternalUserClient
}

// NewUserService constructs a userService injecting the given repository and external client.
func NewUserService(repo UserRepository, externalClient ExternalUserClient) UserService {
	return &userService{
		repo:           repo,
		externalClient: externalClient,
	}
}

func (s *userService) GetAll(ctx context.Context, page int, size int) ([]User, int64, error) {
	users, total, err := s.repo.FindAll(ctx, page, size)
	if err != nil {
		return nil, 0, fmt.Errorf("listing users: %w", err)
	}
	return users, total, nil
}

func (s *userService) GetByID(ctx context.Context, id int64) (*User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if u == nil {
		return nil, ErrNotFound
	}
	return u, nil
}

func (s *userService) Create(ctx context.Context, u User) (User, error) {
	exists, err := s.repo.ExistsByEmail(ctx, u.Email)
	if err != nil {
		return User{}, fmt.Errorf("checking email uniqueness: %w", err)
	}
	if exists {
		return User{}, ErrAlreadyExists
	}

	created, err := s.repo.Save(ctx, u)
	if err != nil {
		return User{}, fmt.Errorf("saving user: %w", err)
	}
	return created, nil
}

func (s *userService) Update(ctx context.Context, id int64, u User) (User, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("finding user for update: %w", err)
	}
	if existing == nil {
		return User{}, ErrNotFound
	}

	u.ID = id
	updated, err := s.repo.Update(ctx, u)
	if err != nil {
		return User{}, fmt.Errorf("updating user: %w", err)
	}
	return updated, nil
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("finding user for delete: %w", err)
	}
	if existing == nil {
		return ErrNotFound
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return nil
}

func (s *userService) FetchFromExternal(ctx context.Context) ([]User, error) {
	users, err := s.externalClient.FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching external users: %w", err)
	}
	return users, nil
}

func (s *userService) SyncFromExternal(ctx context.Context, externalID int64) (User, error) {
	external, err := s.externalClient.FetchByID(ctx, externalID)
	if err != nil {
		return User{}, fmt.Errorf("fetching external user %d: %w", externalID, err)
	}
	if external == nil {
		return User{}, ErrNotFound
	}

	existing, err := s.repo.FindByEmail(ctx, external.Email)
	if err != nil {
		return User{}, fmt.Errorf("looking up local user by email: %w", err)
	}

	if existing != nil {
		external.ID = existing.ID
		updated, err := s.repo.Update(ctx, *external)
		if err != nil {
			return User{}, fmt.Errorf("updating synced user: %w", err)
		}
		return updated, nil
	}

	external.ID = 0
	saved, err := s.repo.Save(ctx, *external)
	if err != nil {
		return User{}, fmt.Errorf("saving synced user: %w", err)
	}
	return saved, nil
}
