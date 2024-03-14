package handlers

import (
	"github.com/go-chi/chi/v5"
	"log"
)

func API(logger *log.Logger) *chi.Mux {
	handler := chi.NewMux()

	handler.Get("/ready", readiness(logger))

	return handler
}
