package middleware

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"runtime/debug"
	"ultimate-service-v1/foundation/web"
)

func Panics(logger *log.Logger) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			values, ok := ctx.Value(web.CtxValues).(*web.Values)
			if !ok {
				return web.NewShutdown("web values missing from context")
			}

			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic: %v", r)

					logger.Printf("%s : PANIC :\n%s", values.TraceId, debug.Stack())
				}
			}()

			return handler(ctx, w, r)
		}
	}
}
