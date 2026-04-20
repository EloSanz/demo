package web

import (
	"context"
	"net/http"
)

type Pinger interface {
	PingContext(ctx context.Context) error
}

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
