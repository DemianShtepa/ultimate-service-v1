package handlers

import (
	"log"
	"net/http"
	"os"
	"ultimate-service-v1/core/middleware"
	"ultimate-service-v1/foundation/web"
)

func API(logger *log.Logger, shutdown chan os.Signal) http.Handler {
	webApp := web.NewApp(shutdown, middleware.Logging(logger), middleware.Errors(logger))

	webApp.Get("/ready", readiness())

	return webApp.Mux
}
