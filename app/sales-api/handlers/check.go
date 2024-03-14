package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func readiness(logger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := struct {
			Status string
		}{
			Status: "OK",
		}

		json.NewEncoder(w).Encode(status)

		logger.Printf("check: readiness called %v", status)
	}
}
