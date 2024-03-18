package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"ultimate-service-v1/foundation/web"
)

func readiness(logger *log.Logger) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		status := struct {
			Status string
		}{
			Status: "OK",
		}

		logger.Printf("check: readiness called %v", status)

		return json.NewEncoder(w).Encode(status)
	}
}
