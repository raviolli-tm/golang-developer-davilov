package broker

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
)

type ConsumerKafka struct {
	consumerCfg *KafkaConf
	Channel     chan kafka.Message
}

func NewConsumer(cfg *KafkaConf) *ConsumerKafka {
	channel := make(chan kafka.Message, 5)
	return &ConsumerKafka{
		consumerCfg: cfg,
		Channel:     channel,
	}
}

func (k *ConsumerKafka) ConsumeWithContext(ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{k.consumerCfg.BrokersUrl},
		Topic:   k.consumerCfg.Topic,
		GroupID: k.consumerCfg.GroupId,
		Logger:  log.New(os.Stdout, "", log.LstdFlags),
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Consumer shutting down...")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v\n", err)
				continue
			}
			k.Channel <- msg
		}
	}

}
