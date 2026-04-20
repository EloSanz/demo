package web

import (
	"log/slog"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

// On error, it delegates to WriteError to produce a consistent JSON error response.
func Adapt(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("handler error", "method", r.Method, "path", r.URL.Path, "error", err)
			WriteError(w, err)
		}
	}
}
