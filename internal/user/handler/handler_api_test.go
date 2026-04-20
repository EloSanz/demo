package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elosanz/demo/internal/user"
	"github.com/elosanz/demo/internal/user/handler"
	"github.com/stretchr/testify/assert"
)

// mockService implements the public user.UserService interface for testing.
type mockService struct {
	getAllFunc            func(ctx context.Context, page, size int) ([]user.User, int64, error)
	getByIDFunc           func(ctx context.Context, id int64) (*user.User, error)
	createFunc            func(ctx context.Context, u user.User) (user.User, error)
	updateFunc            func(ctx context.Context, id int64, u user.User) (user.User, error)
	deleteFunc            func(ctx context.Context, id int64) error
	fetchFromExternalFunc func(ctx context.Context) ([]user.User, error)
	syncFromExternalFunc  func(ctx context.Context, externalID int64) (user.User, error)
	transferFunc          func(ctx context.Context, fromID, toID int64, amount int) error
}

func (m *mockService) GetAll(ctx context.Context, page, size int) ([]user.User, int64, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, page, size)
	}
	return nil, 0, nil
}
func (m *mockService) GetByID(ctx context.Context, id int64) (*user.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockService) Create(ctx context.Context, u user.User) (user.User, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, u)
	}
	return user.User{}, nil
}
func (m *mockService) Update(ctx context.Context, id int64, u user.User) (user.User, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, u)
	}
	return user.User{}, nil
}
func (m *mockService) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}
func (m *mockService) FetchFromExternal(ctx context.Context) ([]user.User, error) {
	if m.fetchFromExternalFunc != nil {
		return m.fetchFromExternalFunc(ctx)
	}
	return nil, nil
}
func (m *mockService) SyncFromExternal(ctx context.Context, externalID int64) (user.User, error) {
	if m.syncFromExternalFunc != nil {
		return m.syncFromExternalFunc(ctx, externalID)
	}
	return user.User{}, nil
}
func (m *mockService) Transfer(ctx context.Context, fromID, toID int64, amount int) error {
	if m.transferFunc != nil {
		return m.transferFunc(ctx, fromID, toID, amount)
	}
	return nil
}

func TestUserHandler_GetAll_Success(t *testing.T) {
	// Arrange
	users := []user.User{{ID: 1, Name: "John", Email: "john@example.com"}}
	mock := &mockService{getAllFunc: func(ctx context.Context, page, size int) ([]user.User, int64, error) {
		return users, int64(len(users)), nil
	}}
	h := handler.NewUserHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/users?page=1&size=10", nil)
	rr := httptest.NewRecorder()

	// Act
	err := h.GetAll(rr, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	var resp struct {
		Content []handler.UserResponse `json:"content"`
		Page    int                    `json:"page"`
		Size    int                    `json:"size"`
		Total   int64                  `json:"total"`
	}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Len(t, resp.Content, 1)
	assert.Equal(t, "John", resp.Content[0].Name)
	assert.Equal(t, int64(1), resp.Total)
}

func TestUserHandler_GetAll_ServiceError(t *testing.T) {
	// Arrange
	mock := &mockService{getAllFunc: func(ctx context.Context, page, size int) ([]user.User, int64, error) {
		return nil, 0, errors.New("db error")
	}}
	h := handler.NewUserHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rr := httptest.NewRecorder()

	// Act
	err := h.GetAll(rr, req)

	// Assert
	assert.Error(t, err)
	// The handler returns a *web.Error (internal) which implements error; we just check status code via the error type if needed.
	// For simplicity, ensure the response code is 500 (internal server error) when the error is written.
	if httpErr, ok := err.(interface{ Status() int }); ok {
		assert.Equal(t, http.StatusInternalServerError, httpErr.Status())
	}
}

// Additional black‑box tests can be added for Create, Update, Delete, etc., following the same pattern.
