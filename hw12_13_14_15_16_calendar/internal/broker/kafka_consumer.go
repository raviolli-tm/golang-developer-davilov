package broker

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type ConsumerKafka struct {
	consumerCfg *KafkaConf
	Channel     chan kafka.Message
}

func NewConsumer(cfg *KafkaConf) *ConsumerKafka {
	return &ConsumerKafka{
		consumerCfg: cfg,
		Channel:     make(chan kafka.Message, 5),
	}
}

// ConsumeWithContext reads messages from the configured topic and forwards
// them to Channel. Channel is closed when the context is cancelled.
func (k *ConsumerKafka) ConsumeWithContext(ctx context.Context) {
	defer close(k.Channel)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.consumerCfg.BrokersUrl},
		Topic:    k.consumerCfg.Topic,
		GroupID:  k.consumerCfg.GroupId,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("failed to close kafka reader: %v", err)
		}
	}()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			log.Printf("error reading kafka message: %v", err)
			continue
		}

		select {
		case k.Channel <- msg:
		case <-ctx.Done():
			return
		}
	}
}

// WaitForBroker blocks until the configured Kafka broker accepts connections
// or the context is cancelled.
func WaitForBroker(ctx context.Context, cfg *KafkaConf) error {
	for {
		conn, err := kafka.DialContext(ctx, "tcp", cfg.BrokersUrl)
		if err == nil {
			_ = conn.Close()
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}
