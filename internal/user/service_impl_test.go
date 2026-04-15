package user_test

import (
	"context"
	"testing"

	"github.com/elosanz/demo/internal/user"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

// ─── Fakes ───────────────────────────────────────────────────────────────────

type fakeUserRepository struct {
	users       map[int64]user.User
	nextID      int64
	emailIndex  map[string]int64
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		users:      make(map[int64]user.User),
		emailIndex: make(map[string]int64),
		nextID:     1,
	}
}

func (f *fakeUserRepository) FindAll(_ context.Context, page int, size int) ([]user.User, int64, error) {
	all := make([]user.User, 0, len(f.users))
	for _, u := range f.users {
		all = append(all, u)
	}
	total := int64(len(all))
	start := (page - 1) * size
	if start >= len(all) {
		return []user.User{}, total, nil
	}
	end := start + size
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], total, nil
}

func (f *fakeUserRepository) FindByID(_ context.Context, id int64) (*user.User, error) {
	u, ok := f.users[id]
	if !ok {
		return nil, nil
	}
	return &u, nil
}

func (f *fakeUserRepository) FindByEmail(_ context.Context, email string) (*user.User, error) {
	id, ok := f.emailIndex[email]
	if !ok {
		return nil, nil
	}
	u := f.users[id]
	return &u, nil
}

func (f *fakeUserRepository) ExistsByEmail(_ context.Context, email string) (bool, error) {
	_, ok := f.emailIndex[email]
	return ok, nil
}

func (f *fakeUserRepository) Save(_ context.Context, u user.User) (user.User, error) {
	u.ID = f.nextID
	f.nextID++
	f.users[u.ID] = u
	f.emailIndex[u.Email] = u.ID
	return u, nil
}

func (f *fakeUserRepository) Update(_ context.Context, u user.User) (user.User, error) {
	if _, ok := f.users[u.ID]; !ok {
		return user.User{}, user.ErrNotFound
	}
	f.users[u.ID] = u
	return u, nil
}

func (f *fakeUserRepository) Delete(_ context.Context, id int64) error {
	u, ok := f.users[id]
	if !ok {
		return user.ErrNotFound
	}
	delete(f.emailIndex, u.Email)
	delete(f.users, id)
	return nil
}

type fakeExternalUserClient struct {
	users  []user.User
	byID   map[int64]*user.User
}

func newFakeExternalUserClient() *fakeExternalUserClient {
	return &fakeExternalUserClient{byID: make(map[int64]*user.User)}
}

func (f *fakeExternalUserClient) FetchAll(_ context.Context) ([]user.User, error) {
	return f.users, nil
}

func (f *fakeExternalUserClient) FetchByID(_ context.Context, id int64) (*user.User, error) {
	return f.byID[id], nil
}

// ─── Factory ─────────────────────────────────────────────────────────────────

type userServiceBuilder struct {
	repo     *fakeUserRepository
	external *fakeExternalUserClient
}

func newUserServiceFactory() *userServiceBuilder {
	return &userServiceBuilder{
		repo:     newFakeUserRepository(),
		external: newFakeExternalUserClient(),
	}
}

func (b *userServiceBuilder) withExistingUser(u user.User) *userServiceBuilder {
	b.repo.users[u.ID] = u
	b.repo.emailIndex[u.Email] = u.ID
	if u.ID >= b.repo.nextID {
		b.repo.nextID = u.ID + 1
	}
	return b
}

func (b *userServiceBuilder) withExternalUser(id int64, u user.User) *userServiceBuilder {
	b.external.byID[id] = &u
	return b
}

func (b *userServiceBuilder) build() user.UserService {
	return user.NewUserService(b.repo, b.external)
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func TestUserService_Create_Success(t *testing.T) {
	// Given
	svc := newUserServiceFactory().build()

	// When
	created, err := svc.Create(ctx, user.User{
		Name:  "Elo",
		Email: "elo@test.com",
		Phone: "123",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "Elo", created.Name)
	require.Equal(t, "elo@test.com", created.Email)
	require.Greater(t, created.ID, int64(0))
}

func TestUserService_Create_Error_WithDuplicateEmail(t *testing.T) {
	// Given
	svc := newUserServiceFactory().
		withExistingUser(user.User{ID: 1, Name: "Existing", Email: "elo@test.com"}).
		build()

	// When
	_, err := svc.Create(ctx, user.User{Name: "New", Email: "elo@test.com"})

	// Then
	require.ErrorIs(t, err, user.ErrAlreadyExists)
}

func TestUserService_GetByID_Success(t *testing.T) {
	// Given
	svc := newUserServiceFactory().
		withExistingUser(user.User{ID: 1, Name: "Elo", Email: "elo@test.com"}).
		build()

	// When
	u, err := svc.GetByID(ctx, 1)

	// Then
	require.NoError(t, err)
	require.NotNil(t, u)
	require.Equal(t, "Elo", u.Name)
}

func TestUserService_GetByID_Error_WithUserNotFound(t *testing.T) {
	// Given
	svc := newUserServiceFactory().build()

	// When
	_, err := svc.GetByID(ctx, 999)

	// Then
	require.ErrorIs(t, err, user.ErrNotFound)
}

func TestUserService_Update_Success(t *testing.T) {
	// Given
	svc := newUserServiceFactory().
		withExistingUser(user.User{ID: 1, Name: "Elo", Email: "elo@test.com"}).
		build()

	// When
	updated, err := svc.Update(ctx, 1, user.User{Name: "Elo Updated", Email: "elo@test.com"})

	// Then
	require.NoError(t, err)
	require.Equal(t, "Elo Updated", updated.Name)
}

func TestUserService_Update_Error_WithUserNotFound(t *testing.T) {
	// Given
	svc := newUserServiceFactory().build()

	// When
	_, err := svc.Update(ctx, 999, user.User{Name: "X", Email: "x@test.com"})

	// Then
	require.ErrorIs(t, err, user.ErrNotFound)
}

func TestUserService_Delete_Success(t *testing.T) {
	// Given
	svc := newUserServiceFactory().
		withExistingUser(user.User{ID: 1, Name: "Elo", Email: "elo@test.com"}).
		build()

	// When
	err := svc.Delete(ctx, 1)

	// Then
	require.NoError(t, err)
}

func TestUserService_Delete_Error_WithUserNotFound(t *testing.T) {
	// Given
	svc := newUserServiceFactory().build()

	// When
	err := svc.Delete(ctx, 999)

	// Then
	require.ErrorIs(t, err, user.ErrNotFound)
}

func TestUserService_SyncFromExternal_Success_WithNewUser(t *testing.T) {
	// Given
	svc := newUserServiceFactory().
		withExternalUser(1, user.User{Name: "External User", Email: "external@test.com", Phone: "456"}).
		build()

	// When
	synced, err := svc.SyncFromExternal(ctx, 1)

	// Then
	require.NoError(t, err)
	require.Equal(t, "External User", synced.Name)
	require.Greater(t, synced.ID, int64(0))
}

func TestUserService_SyncFromExternal_Error_WithExternalUserNotFound(t *testing.T) {
	// Given
	svc := newUserServiceFactory().build()

	// When
	_, err := svc.SyncFromExternal(ctx, 999)

	// Then
	require.ErrorIs(t, err, user.ErrNotFound)
}
