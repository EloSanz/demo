package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elosanz/demo/infrastructure/httpclient"
	"github.com/elosanz/demo/internal/rickandmorty"
	"github.com/elosanz/demo/internal/rickandmorty/handler"
	"github.com/stretchr/testify/require"
)

func setupRMIntegration(_ *testing.T) (*handler.RickAndMortyHandler, *httptest.Server) {
	// 1. Mock Server para Rick and Morty API
	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Mock single character
		if r.URL.Path == "/character/1" {
			json.NewEncoder(w).Encode(map[string]any{
				"id": 1, "name": "Rick Sanchez", "status": "Alive", "species": "Human",
			})
			return
		}

		// Mock search
		if r.URL.Path == "/character/" {
			json.NewEncoder(w).Encode(map[string]any{
				"info": map[string]any{"count": 1, "pages": 1},
				"results": []map[string]any{
					{"id": 1, "name": "Rick Sanchez"},
				},
			})
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))

	// 2. Stack
	client := httpclient.NewRickAndMortyHTTPClient(mockAPI.URL, http.DefaultClient)
	svc := rickandmorty.NewRickAndMortyService(client)

	return handler.NewRickAndMortyHandler(svc), mockAPI
}

func TestRickAndMortyHandler_Integration(t *testing.T) {
	h, server := setupRMIntegration(t)
	defer server.Close()

	t.Run("Get Character by ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/rickandmorty/characters/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		err := h.GetByID(w, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)

		var resp handler.CharacterResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		require.Equal(t, "Rick Sanchez", resp.Name)
	})

	t.Run("Search Characters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/rickandmorty/characters?name=rick", nil)
		w := httptest.NewRecorder()

		err := h.Search(w, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)

		var resp handler.CharacterPageResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		require.Equal(t, 1, resp.Info.Count)
		require.Equal(t, "Rick Sanchez", resp.Results[0].Name)
	})
}
