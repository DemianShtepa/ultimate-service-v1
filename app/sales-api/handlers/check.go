package handlers

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"ultimate-service-v1/foundation/database"
	"ultimate-service-v1/foundation/web"
)

func readiness(db *sqlx.DB) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		statusCode := http.StatusOK
		statusMessage := "OK"

		if err := database.StatusCheck(ctx, db); err != nil {
			fmt.Println(err)
			statusCode = http.StatusInternalServerError
			statusMessage = "db is not ready"
		}

		status := struct {
			Status string
		}{
			Status: statusMessage,
		}

		return web.Respond(ctx, w, status, statusCode)
	}
}
