package user

import (
	"context"
	"fmt"
	"github.com/elosanz/demo/internal/notification"
	"github.com/elosanz/demo/pkg/transaction"
)

// userService implements UserService.
type userService struct {
	repo           UserRepository
	externalClient ExternalUserClient
	notifSvc       notification.NotificationService
}

// NewUserService constructs a userService injecting dependencies.
func NewUserService(repo UserRepository, externalClient ExternalUserClient, notifSvc notification.NotificationService) UserService {
	return &userService{
		repo:           repo,
		externalClient: externalClient,
		notifSvc:       notifSvc,
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
	// 1. Fetch external data
	external, err := s.externalClient.FetchByID(ctx, externalID)
	if err != nil || external == nil {
		return User{}, fmt.Errorf("fetching external user: %w", err)
	}

	tx := transaction.NewHelper(2) // 2 reintentos para el rollback
	var savedUser User

	// 2. Define the Saga Actions
	actions := []transaction.Action{
		{
			Name: "UpsertLocalUser",
			Execute: func() error {
				existing, _ := s.repo.FindByEmail(ctx, external.Email)
				if existing != nil {
					external.ID = existing.ID
				}
				savedUser, err = s.repo.Save(ctx, *external)
				return err
			},
			Rollback: func() error {
				// Si falló el paso siguiente, borramos el usuario creado
				return s.repo.Delete(ctx, savedUser.ID)
			},
		},
		{
			Name: "SendWelcomeNotification",
			Execute: func() error {
				return s.notifSvc.Publish(ctx, notification.Notification{
					Type:    "welcome_email",
					Content: fmt.Sprintf("Welcome %s! Your account is synced.", savedUser.Name),
				})
			},
			// No rollback needed for notification, but could be a cancel-mail if exists
		},
	}

	// 3. Execute Saga
	if err := tx.Execute(actions); err != nil {
		return User{}, fmt.Errorf("sync transaction failed: %w", err)
	}

	return savedUser, nil
}
