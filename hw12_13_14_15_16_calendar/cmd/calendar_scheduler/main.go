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
	flag.StringVar(&configFile, "config", "configs/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Printf("scheduler stopped with error: %v", err)
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

	scannerStorage := storage.NewScannerStorage()
	kafkaProducer := kafka.NewProducer(&config.Broker, &config.RetryPolicy)
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			log.Printf("failed to close kafka producer: %v", err)
		}
	}()

	if _, err := kafkaProducer.GetProducer(ctx); err != nil {
		return fmt.Errorf("connect to kafka: %w", err)
	}

	go func() {
		if err := metrics.Serve(ctx, config.MetricsAddr()); err != nil {
			log.Printf("failed to serve metrics: %v", err)
		}
	}()

	go cleanupOldEvents(ctx, scannerStorage)

	interval := time.Duration(config.Broker.Interval) * time.Second
	if interval <= 0 {
		interval = 5 * time.Second
	}

	cleanupInterval := time.Duration(config.CleanupIntervalHours) * time.Hour
	if cleanupInterval <= 0 {
		cleanupInterval = 24 * time.Hour
	}

	notifyTicker := time.NewTicker(interval)
	cleanupTicker := time.NewTicker(cleanupInterval)
	defer notifyTicker.Stop()
	defer cleanupTicker.Stop()

	lastScan := time.Now()

	for {
		select {
		case <-ctx.Done():
			log.Println("scheduler is shutting down...")
			return nil
		case <-notifyTicker.C:
			now := time.Now()
			sendNotifications(ctx, scannerStorage, kafkaProducer, lastScan, now)
			lastScan = now
		case <-cleanupTicker.C:
			cleanupOldEvents(ctx, scannerStorage)
		}
	}
}

func sendNotifications(
	ctx context.Context,
	scannerStorage *sqlstorage.ScannerStorage,
	producer *kafka.ProducerKafka,
	from, to time.Time,
) {
	metrics.SchedulerScansTotal.Inc()

	scanCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	events, err := scannerStorage.ReadEventToNotify(scanCtx, from, to)
	if err != nil {
		metrics.SchedulerErrorsTotal.Inc()
		log.Printf("failed to read events to notify: %v", err)
		return
	}

	for _, event := range events {
		notification, err := json.Marshal(eventToNotification(event))
		if err != nil {
			metrics.SchedulerErrorsTotal.Inc()
			log.Printf("failed to marshal notification: %v", err)
			continue
		}

		sendCtx, sendCancel := context.WithTimeout(ctx, 10*time.Second)
		err = producer.SendMessage(sendCtx, []byte(event.ID.String()), notification)
		sendCancel()

		if err != nil {
			metrics.SchedulerErrorsTotal.Inc()
			log.Printf("failed to send notification for event %s: %v", event.ID, err)
			continue
		}

		metrics.SchedulerNotificationsSentTotal.Inc()
		log.Printf("notification for event %s has been sent", event.ID)
	}
}

func cleanupOldEvents(ctx context.Context, scannerStorage *sqlstorage.ScannerStorage) {
	cleanupCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	toDelete, err := scannerStorage.ReadEventToDelete(cleanupCtx)
	if err != nil {
		metrics.SchedulerErrorsTotal.Inc()
		log.Printf("failed to read events to delete: %v", err)
		return
	}

	if len(toDelete) == 0 {
		return
	}

	deletedIDs, err := scannerStorage.DeleteEventByIDs(cleanupCtx, toDelete)
	if err != nil {
		metrics.SchedulerErrorsTotal.Inc()
		log.Printf("failed to delete events: %v", err)
		return
	}

	metrics.SchedulerEventsDeletedTotal.Add(float64(len(deletedIDs)))
	log.Printf("deleted %d stale events", len(deletedIDs))
}

func eventToNotification(e storage.Event) storage.Notification {
	return storage.Notification{
		NotificationId:         e.ID,
		NotificationEventTitle: e.Title,
		NotificationEventDate:  e.DateStart,
		NotificationDate:       time.Now(),
		NotificationUserId:     e.UserId,
	}
}
