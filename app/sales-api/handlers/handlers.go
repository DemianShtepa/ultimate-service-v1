package handlers

import (
	"log"
	"net/http"
	"os"
	"ultimate-service-v1/foundation/web"
)

func API(logger *log.Logger, shutdown chan os.Signal) http.Handler {
	webApp := web.NewApp(shutdown)

	webApp.Get("/ready", readiness(logger))

	return webApp.Mux
}
