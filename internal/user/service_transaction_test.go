package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/elosanz/demo/internal/notification"
	"github.com/elosanz/demo/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks necesarios
type mockRepo struct{ mock.Mock }

func (m *mockRepo) FindByID(ctx context.Context, id int64) (*user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}
func (m *mockRepo) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}
func (m *mockRepo) FindAll(ctx context.Context, page, size int) ([]user.User, int64, error) {
	args := m.Called(ctx, page, size)
	return args.Get(0).([]user.User), int64(args.Int(1)), args.Error(2)
}
func (m *mockRepo) Save(ctx context.Context, u user.User) (user.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(user.User), args.Error(1)
}
func (m *mockRepo) Update(ctx context.Context, u user.User) (user.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(user.User), args.Error(1)
}
func (m *mockRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}
func (m *mockRepo) TransferPoints(ctx context.Context, fromID, toID int64, amount int) error {
	args := m.Called(ctx, fromID, toID, amount)
	return args.Error(0)
}

type mockExternal struct{ mock.Mock }

func (m *mockExternal) FetchByID(ctx context.Context, id int64) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}
func (m *mockExternal) FetchAll(ctx context.Context) ([]user.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]user.User), args.Error(1)
}

type mockNotif struct{ mock.Mock }

func (m *mockNotif) Publish(ctx context.Context, n notification.Notification) error {
	args := m.Called(ctx, n)
	return args.Error(0)
}
func (m *mockNotif) StartWorker(ctx context.Context) { m.Called(ctx) }

func TestUserService_SyncFromExternal_RollbackOnNotificationFailure(t *testing.T) {
	repo := new(mockRepo)
	external := new(mockExternal)
	notif := new(mockNotif)
	svc := user.NewUserService(repo, nil, external, notif)

	ctx := context.Background()
	externalUser := &user.User{ID: 0, Name: "Test User", Email: "test@example.com"}
	savedUser := user.User{ID: 123, Name: "Test User", Email: "test@example.com"}

	external.On("FetchByID", ctx, int64(1)).Return(externalUser, nil)
	repo.On("FindByEmail", ctx, externalUser.Email).Return(nil, nil)
	repo.On("Save", ctx, *externalUser).Return(savedUser, nil)
	notif.On("Publish", ctx, mock.Anything).Return(errors.New("notification service down"))
	repo.On("Delete", ctx, int64(123)).Return(nil)

	_, err := svc.SyncFromExternal(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sync transaction failed")

	repo.AssertExpectations(t)
	notif.AssertExpectations(t)
	external.AssertExpectations(t)
}
