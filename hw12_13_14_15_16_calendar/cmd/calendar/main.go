package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/internal/app"
	"github.com/davilov/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/davilov/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/memory"
)

var Logg *logger.Logger
var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	Logg = logger.New(config.Logger.Level)

	storage := memorystorage.New()
	calendar := app.New(Logg, storage)

	server := internalhttp.NewServer(Logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			Logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	Logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		Logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
