package main

import (
	"context"
	"encoding/json"
	"flag"
	kafka "github.com/davilov/hw12_13_14_15_calendar/internal/broker"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/davilov/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/google/uuid"
	"log"
	"os/signal"
	"syscall"
	"time"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	defer cancel()

	config := NewConfig()

	scannerStorage := sqlstorage.New(config.Database, ctx).NewScannerStorage()
	kafkaProducer := kafka.NewProducer(&config.Broker, &config.RetryPolicy)

	_, err := kafkaProducer.GetProducer(ctx)
	if err != nil {
		log.Printf("Failed to detect broker: %v", err)
		return
	}

	tickerNotify := time.NewTicker(time.Duration(config.Broker.Interval) * time.Second)
	tickerDelete := time.NewTicker(time.Hour * 24)
	defer tickerNotify.Stop()
	defer tickerDelete.Stop()

	eventTime := time.Now()

	for {
		select {

		case <-ctx.Done():
			notification, _ := json.Marshal(&storage.Notification{
				NotificationId:         uuid.New(),
				NotificationEventTitle: "",
				NotificationEventDate:  time.Now(),
				NotificationDate:       time.Now().Add(10 * time.Minute),
				NotificationUserId:     0,
			})
			err = kafkaProducer.SendMessage(
				context.Background(),
				[]byte(time.Now().Format("2006-01-02-15-04-05")),
				notification)

			log.Println("Shutting down producer...")
			kafkaProducer.Close()
			return

		case <-tickerNotify.C:
			go func() {
				ctxScan, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()
				events, err := scannerStorage.ReadEventToNotify(ctxScan, eventTime.Add(time.Minute), config.Broker.Interval)
				if err != nil {
					log.Printf("Failed to read event to notify: %v", err)
					return
				}

				for _, e := range events {
					notification, err := json.Marshal(EventToNotification(e))
					if err != nil {
						log.Printf("Failed to marshal notification: %v", err)
						continue
					}

					err = kafkaProducer.SendMessage(
						context.Background(),
						[]byte(time.Now().Format("2006-01-02-15-04-05")),
						notification)
					if err != nil {
						log.Printf("Failed to send notification: %v", err)
						continue
					}
				}

			}()
		case <-tickerDelete.C:
			go func() {
				ctxScan, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				toDelete, err := scannerStorage.ReadEventToDelete(ctxScan)
				defer cancel()
				if err != nil {
					log.Printf("Failed to read event to delete: %v", err)
					return
				}
				deletedIDs, err := scannerStorage.DeleteEventByIDs(ctxScan, toDelete)
				if err != nil {
					log.Printf("Failed to delete events: %v", err)
					return
				}
				log.Printf("Deleted %d events", len(deletedIDs))
			}()
		}
	}
}

func EventToNotification(e storage.Event) storage.Notification {
	n := storage.Notification{}
	n.NotificationId = e.ID
	n.NotificationUserId = e.UserId
	n.NotificationEventTitle = e.Title
	n.NotificationEventDate = e.DateStart
	return n
}
