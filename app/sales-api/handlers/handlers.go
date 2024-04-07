package handlers

import (
	"log"
	"net/http"
	"os"
	"ultimate-service-v1/core/authentication"
	"ultimate-service-v1/core/middleware"
	"ultimate-service-v1/foundation/web"
)

func API(logger *log.Logger, shutdown chan os.Signal, auth *authentication.Authentication) http.Handler {
	webApp := web.NewApp(
		shutdown,
		middleware.Logging(logger),
		middleware.Errors(logger),
		middleware.Metrics(),
		middleware.Panics(logger),
	)

	webApp.Get("/ready", readiness(), middleware.Authenticate(auth), middleware.Authorize(auth, authentication.RoleAdmin))

	return webApp.Mux
}
