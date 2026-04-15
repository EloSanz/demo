// Package handler implements the HTTP transport layer for the storage domain.
package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/elosanz/demo/internal/storage"
	"github.com/elosanz/demo/pkg/web"
)

// StorageHandler handles HTTP requests for the storage domain.
type StorageHandler struct {
	service storage.StorageService
}

// NewStorageHandler constructs a StorageHandler with the given service.
func NewStorageHandler(service storage.StorageService) *StorageHandler {
	return &StorageHandler{service: service}
}

// Upload handles POST /api/storage/upload (multipart/form-data, field: file)
func (h *StorageHandler) Upload(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32 MB max
		return web.NewError(http.StatusBadRequest, "invalid multipart form")
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return web.NewError(http.StatusBadRequest, "file field is required")
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed reading uploaded file", "error", err)
		return web.NewError(http.StatusInternalServerError, "internal server error")
	}

	message, err := h.service.Upload(r.Context(), header.Filename, content)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed uploading file to S3", "filename", header.Filename, "error", err)
		return web.NewError(http.StatusInternalServerError, "internal server error")
	}

	return web.EncodeJSON(w, map[string]string{"message": message}, http.StatusOK)
}

// List handles GET /api/storage/files
func (h *StorageHandler) List(w http.ResponseWriter, r *http.Request) error {
	files, err := h.service.List(r.Context())
	if err != nil {
		slog.ErrorContext(r.Context(), "failed listing S3 files", "error", err)
		return web.NewError(http.StatusInternalServerError, "internal server error")
	}
	return web.EncodeJSON(w, files, http.StatusOK)
}

// Download handles GET /api/storage/files/{name}
func (h *StorageHandler) Download(w http.ResponseWriter, r *http.Request) error {
	name, err := web.Param(r, "name")
	if err != nil {
		return err
	}

	content, err := h.service.Download(r.Context(), name)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			return web.NewError(http.StatusNotFound, err.Error())
		default:
			slog.ErrorContext(r.Context(), "failed downloading file from S3", "name", name, "error", err)
			return web.NewError(http.StatusInternalServerError, "internal server error")
		}
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
	return nil
}

// Delete handles DELETE /api/storage/files/{name}
func (h *StorageHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	name, err := web.Param(r, "name")
	if err != nil {
		return err
	}

	if err := h.service.Delete(r.Context(), name); err != nil {
		slog.ErrorContext(r.Context(), "failed deleting file from S3", "name", name, "error", err)
		return web.NewError(http.StatusInternalServerError, "internal server error")
	}

	return web.EncodeJSON(w, nil, http.StatusNoContent)
}
