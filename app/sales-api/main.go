package main

import (
	"context"
	"errors"
	_ "expvar"
	"fmt"
	"github.com/ardanlabs/conf/v3"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ultimate-service-v1/app/sales-api/handlers"
)

var build = "develop"

func main() {
	logger := log.New(os.Stdout, "SALES: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(logger); err != nil {
		logger.Println("main: error", err)
	}
}

func run(logger *log.Logger) error {
	config := struct {
		conf.Version
		Web struct {
			ApiHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
	}{}
	config.Version.Build = build
	config.Version.Desc = "Project description."

	info, err := conf.Parse("SALES", &config)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(info)

			return nil
		}

		return err
	}

	log.Printf("run: Application started. Version %q", build)
	defer log.Println("run: Completed")

	out, err := conf.String(&config)
	if err != nil {
		return err
	}
	log.Printf("run: Config: \n%v\n", out)

	log.Println("run: Initialize debugging support")
	go func() {
		log.Printf("run: Debug listening %s", config.Web.DebugHost)
		if err := http.ListenAndServe(config.Web.DebugHost, http.DefaultServeMux); err != nil {
			log.Printf("run: Debug listener closed : %v", err)
		}
	}()

	log.Println("run: Initialize api server")
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	api := http.Server{
		Addr:         config.Web.ApiHost,
		Handler:      handlers.API(),
		ReadTimeout:  config.Web.ReadTimeout,
		WriteTimeout: config.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("run: Starting api server %s", config.Web.ApiHost)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return err
	case sig := <-shutdown:
		log.Printf("run: Starting shutdown - %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), config.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()

			return err
		}
	}

	return nil
}
