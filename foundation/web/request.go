package web

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"net/http"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func Decode(r *http.Request, value any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(value); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		fieldErrors := make([]FieldError, len(validationErrors))
		for i, validationError := range validationErrors {
			fieldErrors[i] = FieldError{
				Field: validationError.Field(),
				Error: validationError.Error(),
			}
		}

		return &Error{
			Value:  errors.New("validation error"),
			Status: http.StatusBadRequest,
			Fields: fieldErrors,
		}
	}

	return nil
}
