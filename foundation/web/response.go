package web

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	values, ok := ctx.Value(CtxValues).(*Values)
	if !ok {
		return NewShutdown("web values missing from context")
	}

	values.StatusCode = statusCode

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)

		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func RespondError(ctx context.Context, w http.ResponseWriter, err error) error {
	if rootErr, ok := errors.Cause(err).(*Error); ok {
		errorResponse := ErrorResponse{
			Error:  rootErr.Error(),
			Fields: rootErr.Fields,
		}

		if err = Respond(ctx, w, errorResponse, rootErr.Status); err != nil {
			return err
		}

		return nil
	}

	errorResponse := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	if err = Respond(ctx, w, errorResponse, http.StatusInternalServerError); err != nil {
		return err
	}

	return nil
}
