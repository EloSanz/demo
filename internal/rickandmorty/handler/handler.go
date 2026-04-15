// Package handler implements the HTTP transport layer for the Rick and Morty domain.
package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/elosanz/demo/internal/rickandmorty"
	"github.com/elosanz/demo/pkg/web"
)

// RickAndMortyHandler handles HTTP requests for the Rick and Morty domain.
type RickAndMortyHandler struct {
	service rickandmorty.RickAndMortyService
}

// NewRickAndMortyHandler constructs a RickAndMortyHandler with the given service.
func NewRickAndMortyHandler(service rickandmorty.RickAndMortyService) *RickAndMortyHandler {
	return &RickAndMortyHandler{service: service}
}

// GetByID handles GET /api/rickandmorty/characters/{id}
func (h *RickAndMortyHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	id, err := web.ParamInt(r, "id")
	if err != nil {
		return err
	}

	character, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, rickandmorty.ErrCharacterNotFound):
			return web.NewError(http.StatusNotFound, err.Error())
		default:
			slog.ErrorContext(r.Context(), "unexpected error fetching character", "id", id, "error", err)
			return web.NewError(http.StatusInternalServerError, "internal server error")
		}
	}

	return web.EncodeJSON(w, mapCharacterToResponse(*character), http.StatusOK)
}

// Search handles GET /api/rickandmorty/characters?page=1&name=rick&status=alive&species=human
func (h *RickAndMortyHandler) Search(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	page := queryIntParam(query.Get("page"), 0)
	name := query.Get("name")
	status := query.Get("status")
	species := query.Get("species")

	result, err := h.service.GetCharacters(r.Context(), page, name, status, species)
	if err != nil {
		slog.ErrorContext(r.Context(), "unexpected error searching characters", "error", err)
		return web.NewError(http.StatusInternalServerError, "internal server error")
	}

	return web.EncodeJSON(w, mapPageToResponse(result), http.StatusOK)
}

// queryIntParam parses a query param string to int; returns defaultVal when empty or invalid.
func queryIntParam(raw string, defaultVal int) int {
	if raw == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return defaultVal
	}
	return v
}
