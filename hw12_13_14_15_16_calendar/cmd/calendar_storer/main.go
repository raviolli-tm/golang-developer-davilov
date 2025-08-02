package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/sql"
	"log"
	"os/signal"
	"syscall"
	"time"

	kafka "github.com/davilov/hw12_13_14_15_calendar/internal/broker"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/storer_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	defer cancel()

	config := NewConfig()
	notificationStorage := sqlstorage.New(config.Database, ctx).NewNotificationStorage()
	consumer := kafka.NewConsumer(&config.Broker)

	go consumer.ConsumeWithContext(ctx)

	for {

		select {
		case <-ctx.Done():
			return

		case notification := <-consumer.Channel:
			ctxScan, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			var n storage.Notification
			err := json.Unmarshal(notification.Value, &n)
			if err != nil {
				log.Printf("Failed to unmarshal notification: %v", err)
				break
			}

			createdNotification, err := notificationStorage.CreateNotification(ctxScan, n)
			if err != nil {
				log.Printf("Failed to create notification: %v", err)
				break
			}
			log.Printf("Created notification: %v", createdNotification)
			cancel()

		}
	}

}
