package main

import (
	"errors"
	_ "expvar"
	"fmt"
	"github.com/ardanlabs/conf/v3"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
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

	select {}

	return nil
}
