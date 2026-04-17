package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elosanz/demo/infrastructure/httpclient"
	"github.com/elosanz/demo/infrastructure/postgres"
	"github.com/elosanz/demo/internal/user"
	"github.com/elosanz/demo/internal/user/handler"
	"github.com/stretchr/testify/require"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupUserIntegration(t *testing.T) (*handler.UserHandler, *gorm.DB, *httptest.Server) {
	mockExternal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/users/99" {
			json.NewEncoder(w).Encode(map[string]any{
				"id": 99, "name": "External User", "email": "ext@test.com",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	db.AutoMigrate(&user.User{})

	repo := postgres.NewUserGORMRepository(db)
	client := httpclient.NewJSONPlaceholderClient(mockExternal.URL, http.DefaultClient)
	svc := user.NewUserService(repo, nil, client, nil)

	return handler.NewUserHandler(svc), db, mockExternal
}

func TestUserHandler_Integration_CompleteFlow(t *testing.T) {
	t.Run("Create and Get workflow", func(t *testing.T) {
		h, db, server := setupUserIntegration(t)
		defer server.Close()
		sqlDB, _ := db.DB()
		defer sqlDB.Close()

		// 1. Create
		body, _ := json.Marshal(handler.CreateUserRequest{Name: "Elo", Email: "elo@test.com"})
		req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		require.NoError(t, h.Create(w, req))
		require.Equal(t, http.StatusCreated, w.Code)

		// 2. Get
		reqGet := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
		reqGet.SetPathValue("id", "1")
		wGet := httptest.NewRecorder()
		require.NoError(t, h.GetByID(wGet, reqGet))
		require.Equal(t, http.StatusOK, wGet.Code)
	})

	t.Run("Sync External User", func(t *testing.T) {
		h, db, server := setupUserIntegration(t)
		defer server.Close()
		sqlDB, _ := db.DB()
		defer sqlDB.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/users/sync/99", nil)
		req.SetPathValue("externalID", "99")
		w := httptest.NewRecorder()
		require.NoError(t, h.SyncExternal(w, req))
		require.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Update and Delete workflow", func(t *testing.T) {
		h, db, server := setupUserIntegration(t)
		defer server.Close()
		sqlDB, _ := db.DB()
		defer sqlDB.Close()

		// Create first
		db.Create(&user.User{ID: 1, Name: "Old", Email: "old@test.com"})

		// Update
		body, _ := json.Marshal(handler.UpdateUserRequest{Name: "New", Email: "old@test.com"})
		reqUpd := httptest.NewRequest(http.MethodPut, "/api/users/1", bytes.NewBuffer(body))
		reqUpd.SetPathValue("id", "1")
		require.NoError(t, h.Update(httptest.NewRecorder(), reqUpd))

		// Delete
		reqDel := httptest.NewRequest(http.MethodDelete, "/api/users/1", nil)
		reqDel.SetPathValue("id", "1")
		wDel := httptest.NewRecorder()
		require.NoError(t, h.Delete(wDel, reqDel))
		require.Equal(t, http.StatusNoContent, wDel.Code)
	})
}
