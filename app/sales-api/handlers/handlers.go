package handlers

import (
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"ultimate-service-v1/core/authentication"
	"ultimate-service-v1/core/middleware"
	"ultimate-service-v1/foundation/web"
)

func API(logger *log.Logger, shutdown chan os.Signal, auth *authentication.Authentication, db *sqlx.DB) http.Handler {
	webApp := web.NewApp(
		shutdown,
		middleware.Logging(logger),
		middleware.Errors(logger),
		middleware.Metrics(),
		middleware.Panics(logger),
	)

	webApp.Get("/ready", readiness(db), middleware.Authenticate(auth), middleware.Authorize(auth, authentication.RoleAdmin))

	return webApp.Mux
}
