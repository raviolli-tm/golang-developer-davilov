package main

import (
	"context"
	"flag"
	memorystorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/internal/app"
	"github.com/davilov/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/davilov/hw12_13_14_15_calendar/internal/server/http"
)

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
	logg := logger.New(config.Logger)
	var storage app.Storage

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	switch config.StorageType {
	case "postgres":
		storage = sqlstorage.New(config.Database, ctx)
	case "in-memory":
		storage = memorystorage.New()
	default:
		storage = memorystorage.New()
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
