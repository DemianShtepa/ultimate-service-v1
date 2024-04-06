package middleware

import (
	"context"
	"log"
	"net/http"
	"ultimate-service-v1/foundation/web"
)

func Errors(logger *log.Logger) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			values, ok := ctx.Value(web.CtxValues).(*web.Values)
			if !ok {
				return web.NewShutdown("web values missing from context")
			}

			if err := handler(ctx, w, r); err != nil {
				logger.Printf("%s : ERROR : %v", values.TraceId, err)

				if err := web.RespondError(ctx, w, err); err != nil {
					return err
				}

				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}
	}
}
