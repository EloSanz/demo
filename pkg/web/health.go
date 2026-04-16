package web

import (
	"context"
	"net/http"
)

// Pinger is an interface for components that can be checked for health.
// sql.DB implements PingContext.
type Pinger interface {
	PingContext(ctx context.Context) error
}

// HealthHandler returns a simple 200 OK if the application is alive.
// Optionally checks a list of Pingers (database, cache, etc.).
func HealthHandler(pingers ...Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, p := range pingers {
			if err := p.PingContext(r.Context()); err != nil {
				WriteError(w, NewError(http.StatusServiceUnavailable, "database is down"))
				return
			}
		}

		EncodeJSON(w, map[string]string{"status": "UP"}, http.StatusOK)
	}
}
