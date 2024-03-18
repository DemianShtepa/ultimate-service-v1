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

const KeyValues ctxKey = 1

type Values struct {
	TraceId string
	Now     time.Time
}

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	Mux      *http.ServeMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		Mux:      http.NewServeMux(),
		shutdown: shutdown,
	}
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) Get(path string, handler Handler) {
	a.handle("GET", path, handler)
}

func (a *App) handle(method string, path string, handler Handler) {
	a.Mux.HandleFunc(fmt.Sprintf("%s %s", method, path), func(w http.ResponseWriter, r *http.Request) {
		values := Values{
			TraceId: uuid.New().String(),
			Now:     time.Time{},
		}
		ctx := context.WithValue(r.Context(), KeyValues, &values)
		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}
	})
}
