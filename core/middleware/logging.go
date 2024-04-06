package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
	"ultimate-service-v1/foundation/web"
)

func Logging(logger *log.Logger) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			values, ok := ctx.Value(web.CtxValues).(*web.Values)
			if !ok {
				return web.NewShutdown("web values missing from context")
			}

			logger.Printf("%s : started   : %s %s -> %s", values.TraceId, r.Method, r.URL.Path, r.RemoteAddr)

			err := handler(ctx, w, r)

			logger.Printf(
				"%s : completed : %s %s -> %s (%d) (%s)",
				values.TraceId,
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				values.StatusCode,
				time.Since(values.Now),
			)

			return err
		}
	}
}
