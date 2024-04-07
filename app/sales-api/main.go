package main

import (
	"context"
	"crypto/rsa"
	"errors"
	_ "expvar"
	"fmt"
	"github.com/ardanlabs/conf/v3"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ultimate-service-v1/app/sales-api/handlers"
	"ultimate-service-v1/core/authentication"
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
		Auth struct {
			KeyId   string `conf:"default:1b372655-7489-4417-b972-c9b58ad6899d"`
			KeyPath string `conf:"default:/Users/demian/GolandProjects/ultimate-service-v1/private.pem"`
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

	logger.Printf("run: Application started. Version %q", build)
	defer logger.Println("run: Completed")

	out, err := conf.String(&config)
	if err != nil {
		return err
	}
	logger.Printf("run: Config: \n%v\n", out)

	logger.Println("run: Initialize debugging support")
	go func() {
		logger.Printf("run: Debug listening %s", config.Web.DebugHost)
		if err := http.ListenAndServe(config.Web.DebugHost, http.DefaultServeMux); err != nil {
			logger.Printf("run: Debug listener closed : %v", err)
		}
	}()

	logger.Println("run: Initialize api server")
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	privateKeyFile, err := os.ReadFile(config.Auth.KeyPath)
	if err != nil {
		return err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		return err
	}
	auth := authentication.NewAuthentication(map[string]*rsa.PrivateKey{
		config.Auth.KeyId: privateKey,
	})

	api := http.Server{
		Addr:         config.Web.ApiHost,
		Handler:      handlers.API(logger, shutdown, auth),
		ReadTimeout:  config.Web.ReadTimeout,
		WriteTimeout: config.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Printf("run: Starting api server %s", config.Web.ApiHost)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return err
	case sig := <-shutdown:
		logger.Printf("run: Starting shutdown - %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), config.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()

			return err
		}
	}

	return nil
}
