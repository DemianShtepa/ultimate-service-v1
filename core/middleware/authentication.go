package middleware

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"ultimate-service-v1/core/authentication"
	"ultimate-service-v1/foundation/web"
)

var ErrForbidden = web.NewRequestError(errors.New("you are not authorized for that action"), http.StatusForbidden)

func Authenticate(a *authentication.Authentication) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			authHeader := r.Header.Get("Authorization")
			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return web.NewRequestError(errors.New("expected bearer authorization header"), http.StatusBadRequest)
			}

			claims, err := a.Authenticate(parts[1])
			if err != nil {
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			ctx = context.WithValue(ctx, authentication.CtxValues, claims)

			return handler(ctx, w, r)
		}
	}
}

func Authorize(a *authentication.Authentication, roles ...string) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, ok := ctx.Value(authentication.CtxValues).(authentication.Claims)
			if !ok {
				return errors.New("claims is missed from context")
			}

			if !a.Authorize(claims, roles...) {
				return ErrForbidden
			}

			return handler(ctx, w, r)
		}
	}
}
