package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	kafka "github.com/davilov/hw12_13_14_15_calendar/internal/broker"
	"github.com/davilov/hw12_13_14_15_calendar/internal/metrics"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/storer_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Printf("storer stopped with error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := NewConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := sqlstorage.New(config.Database, ctx)
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer closeCancel()
		_ = storage.Close(closeCtx)
	}()

	notificationStorage := storage.NewNotificationStorage()

	if err := kafka.WaitForBroker(ctx, &config.Broker); err != nil {
		return fmt.Errorf("wait for kafka: %w", err)
	}

	consumer := kafka.NewConsumer(&config.Broker)
	go consumer.ConsumeWithContext(ctx)

	go func() {
		if err := metrics.Serve(ctx, config.MetricsAddr()); err != nil {
			log.Printf("failed to serve metrics: %v", err)
		}
	}()

	log.Println("storer is running...")

	for {
		select {
		case <-ctx.Done():
			log.Println("storer is shutting down...")
			return nil
		case message, ok := <-consumer.Channel:
			if !ok {
				return nil
			}

			processNotification(ctx, notificationStorage, message.Value)
		}
	}
}

func processNotification(ctx context.Context, store *sqlstorage.NotificationStorage, payload []byte) {
	var notification storage.Notification
	if err := json.Unmarshal(payload, &notification); err != nil {
		metrics.NotificationsErrorsTotal.Inc()
		log.Printf("failed to unmarshal notification: %v", err)
		return
	}

	insertCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	created, err := store.CreateNotification(insertCtx, notification)
	if err != nil {
		metrics.NotificationsErrorsTotal.Inc()
		log.Printf("failed to create notification: %v", err)
		return
	}

	metrics.NotificationsStoredTotal.Inc()
	log.Printf("notification %s has been stored", created.NotificationId)
}
