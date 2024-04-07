package web

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"syscall"
	"time"
)

type ctxKey int

const CtxValues ctxKey = 1

type Values struct {
	TraceId    string
	Now        time.Time
	StatusCode int
}

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	Mux         *http.ServeMux
	shutdown    chan os.Signal
	middlewares []Middleware
}

func NewApp(shutdown chan os.Signal, middlewares ...Middleware) *App {
	return &App{
		Mux:         http.NewServeMux(),
		shutdown:    shutdown,
		middlewares: middlewares,
	}
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) Get(path string, handler Handler, middlewares ...Middleware) {
	a.handle("GET", path, handler, middlewares...)
}

func (a *App) handle(method string, path string, handler Handler, middlewares ...Middleware) {
	handler = wrapMiddleware(middlewares, handler)
	handler = wrapMiddleware(a.middlewares, handler)

	a.Mux.HandleFunc(fmt.Sprintf("%s %s", method, path), func(w http.ResponseWriter, r *http.Request) {
		values := Values{
			TraceId: uuid.New().String(),
			Now:     time.Now(),
		}
		ctx := context.WithValue(r.Context(), CtxValues, &values)
		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	})
}
