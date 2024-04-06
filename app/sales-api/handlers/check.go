package handlers

import (
	"context"
	"net/http"
	"ultimate-service-v1/foundation/web"
)

func readiness() web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		status := struct {
			Status string
		}{
			Status: "OK",
		}

		return web.Respond(ctx, w, status, http.StatusOK)
	}
}
