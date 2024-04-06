package middleware

import (
	"context"
	"expvar"
	"net/http"
	"runtime"
	"ultimate-service-v1/foundation/web"
)

var metrics = struct {
	requests  *expvar.Int
	gorutines *expvar.Int
	errors    *expvar.Int
}{
	requests:  expvar.NewInt("requests"),
	gorutines: expvar.NewInt("gorutines"),
	errors:    expvar.NewInt("errors"),
}

func Metrics() web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			metrics.requests.Add(1)

			if metrics.requests.Value()%100 == 0 {
				metrics.gorutines.Set(int64(runtime.NumGoroutine()))
			}

			err := handler(ctx, w, r)
			if err != nil {
				metrics.errors.Add(1)
			}

			return err
		}
	}
}
