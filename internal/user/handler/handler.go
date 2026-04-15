// Package handler implements the HTTP transport layer for the user domain.
package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/elosanz/demo/internal/user"
	"github.com/elosanz/demo/pkg/web"
)

// UserHandler handles HTTP requests for the user domain.
type UserHandler struct {
	service user.UserService
}

// NewUserHandler constructs a UserHandler with the given service.
func NewUserHandler(service user.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetAll handles GET /api/users?page=1&size=10
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	page := queryInt(r, "page", 1)
	size := queryInt(r, "size", 10)

	users, total, err := h.service.GetAll(r.Context(), page, size)
	if err != nil {
		slog.ErrorContext(r.Context(), "unexpected error listing users", "error", err)
		return web.NewError(http.StatusInternalServerError, "internal server error")
	}

	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = mapToResponse(u)
	}

	return web.EncodeJSON(w, PageResponse{
		Content: responses,
		Page:    page,
		Size:    size,
		Total:   total,
	}, http.StatusOK)
}

// GetByID handles GET /api/users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	id, err := web.ParamInt(r, "id")
	if err != nil {
		return err
	}

	u, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return web.NewError(http.StatusNotFound, err.Error())
		default:
			slog.ErrorContext(r.Context(), "unexpected error fetching user", "id", id, "error", err)
			return web.NewError(http.StatusInternalServerError, "internal server error")
		}
	}

	response := mapToResponse(*u)
	return web.EncodeJSON(w, response, http.StatusOK)
}

// Create handles POST /api/users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var req CreateUserRequest
	if err := web.DecodeJSON(r, &req); err != nil {
		return err
	}

	if err := validateCreateRequest(req); err != nil {
		return web.NewError(http.StatusUnprocessableEntity, err.Error())
	}

	created, err := h.service.Create(r.Context(), req.toDomain())
	if err != nil {
		switch {
		case errors.Is(err, user.ErrAlreadyExists):
			return web.NewError(http.StatusConflict, err.Error())
		default:
			slog.ErrorContext(r.Context(), "unexpected error creating user", "error", err)
			return web.NewError(http.StatusInternalServerError, "internal server error")
		}
	}

	return web.EncodeJSON(w, mapToResponse(created), http.StatusCreated)
}

// Update handles PUT /api/users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) error {
	id, err := web.ParamInt(r, "id")
	if err != nil {
		return err
	}

	var req UpdateUserRequest
	if err := web.DecodeJSON(r, &req); err != nil {
		return err
	}

	if err := validateUpdateRequest(req); err != nil {
		return web.NewError(http.StatusUnprocessableEntity, err.Error())
	}

	updated, err := h.service.Update(r.Context(), id, req.toDomain())
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return web.NewError(http.StatusNotFound, err.Error())
		default:
			slog.ErrorContext(r.Context(), "unexpected error updating user", "id", id, "error", err)
			return web.NewError(http.StatusInternalServerError, "internal server error")
		}
	}

	return web.EncodeJSON(w, mapToResponse(updated), http.StatusOK)
}

// Delete handles DELETE /api/users/{id}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id, err := web.ParamInt(r, "id")
	if err != nil {
		return err
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return web.NewError(http.StatusNotFound, err.Error())
		default:
			slog.ErrorContext(r.Context(), "unexpected error deleting user", "id", id, "error", err)
			return web.NewError(http.StatusInternalServerError, "internal server error")
		}
	}

	return web.EncodeJSON(w, nil, http.StatusNoContent)
}

// FetchExternal handles GET /api/users/external
func (h *UserHandler) FetchExternal(w http.ResponseWriter, r *http.Request) error {
	users, err := h.service.FetchFromExternal(r.Context())
	if err != nil {
		slog.ErrorContext(r.Context(), "unexpected error fetching external users", "error", err)
		return web.NewError(http.StatusInternalServerError, "internal server error")
	}

	responses := make([]UserResponse, len(users))
	for i, u := range users {
		responses[i] = mapToResponse(u)
	}

	return web.EncodeJSON(w, responses, http.StatusOK)
}

// SyncExternal handles POST /api/users/sync/{externalID}
func (h *UserHandler) SyncExternal(w http.ResponseWriter, r *http.Request) error {
	externalID, err := web.ParamInt(r, "externalID")
	if err != nil {
		return err
	}

	synced, err := h.service.SyncFromExternal(r.Context(), externalID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return web.NewError(http.StatusNotFound, err.Error())
		default:
			slog.ErrorContext(r.Context(), "unexpected error syncing external user", "externalID", externalID, "error", err)
			return web.NewError(http.StatusInternalServerError, "internal server error")
		}
	}

	return web.EncodeJSON(w, mapToResponse(synced), http.StatusCreated)
}

// queryInt extracts an integer query parameter, returning defaultVal if absent or invalid.
func queryInt(r *http.Request, name string, defaultVal int) int {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < 1 {
		return defaultVal
	}
	return v
}

// validateCreateRequest checks required fields on a create request.
func validateCreateRequest(req CreateUserRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

// validateUpdateRequest checks required fields on an update request.
func validateUpdateRequest(req UpdateUserRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	return nil
}
