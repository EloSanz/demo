package web

import (
	"log/slog"
	"net/http"
)

// HandlerFunc is an HTTP handler that returns an error.
// Handlers return nil on success; return web.NewError for 4xx/5xx responses.
// The framework (Adapt) is responsible for writing the error response.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// Adapt converts a HandlerFunc into a standard http.HandlerFunc.
// On error, it delegates to WriteError to produce a consistent JSON error response.
func Adapt(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("handler error", "method", r.Method, "path", r.URL.Path, "error", err)
			WriteError(w, err)
		}
	}
}
