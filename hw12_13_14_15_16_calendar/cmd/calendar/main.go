package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/internal/app"
	"github.com/davilov/hw12_13_14_15_calendar/internal/logger"
	"github.com/davilov/hw12_13_14_15_calendar/internal/metrics"
	internalhttp "github.com/davilov/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/sql"
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

	if err := run(); err != nil {
		log.Printf("calendar stopped with error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := NewConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logg := logger.New(config.Logger)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := newStorage(config, ctx)
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer closeCancel()

		if err := storage.Close(closeCtx); err != nil {
			logg.Error("failed to close storage: " + err.Error())
		}
	}()

	go func() {
		if err := metrics.Serve(ctx, config.MetricsAddr()); err != nil && !errors.Is(err, context.Canceled) {
			logg.Error("failed to serve metrics: " + err.Error())
		}
	}()

	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar, config.HTTPAddr())

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}

	return nil
}

func newStorage(config Config, ctx context.Context) (app.Storage, error) {
	switch config.StorageType {
	case "postgres":
		storage, err := sqlstorage.New(config.Database, ctx)
		if err != nil {
			return nil, fmt.Errorf("connect to postgres: %w", err)
		}

		return storage.NewEventStorage(), nil
	case "in-memory":
		return memorystorage.New(), nil
	default:
		return nil, fmt.Errorf("unknown storage type %q", config.StorageType)
	}
}
