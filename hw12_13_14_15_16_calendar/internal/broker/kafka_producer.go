package broker

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

type ProducerKafka struct {
	producerCfg *KafkaConf
	retryCfg    *RetryPolicyConf
	producer    *kafka.Writer
	mu          sync.Mutex
}

func NewProducer(cfg *KafkaConf, retryCfg *RetryPolicyConf) *ProducerKafka {
	return &ProducerKafka{
		producerCfg: cfg,
		retryCfg:    retryCfg,
		mu:          sync.Mutex{},
	}
}

func (k *ProducerKafka) GetProducer(ctx context.Context) (*kafka.Writer, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.producer != nil {
		return k.producer, nil
	}
	var finalErr error
	for i := 0; i < k.retryCfg.MaxRetries; i++ {

		fmt.Printf("Attempt %d to create producer\n", i+1)
		conn, finalErr := kafka.DialContext(ctx, "tcp", k.producerCfg.BrokersUrl)
		if finalErr != nil {
			finalErr = fmt.Errorf("ошибка подключения к брокеру: %w", finalErr)
			if i+1 != k.retryCfg.MaxRetries {
				time.Sleep(time.Second * time.Duration(k.retryCfg.BaseTime))
			}
			continue
		}
		err := conn.Close()
		if err != nil {
			return nil, err
		}

		k.producer = &kafka.Writer{
			Addr:        kafka.TCP(k.producerCfg.BrokersUrl),
			Topic:       k.producerCfg.Topic,
			Balancer:    &kafka.LeastBytes{},
			MaxAttempts: 3,
		}

		return k.producer, nil

	}
	return k.producer, finalErr

}

func (k *ProducerKafka) SendMessage(ctx context.Context, key []byte, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	writer, err := k.GetProducer(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to broker: %v", err)
	}
	err = writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (k *ProducerKafka) Close() {
	k.mu.Lock()
	defer k.mu.Unlock()
	_ = k.producer.Close()
	k.producer = nil

}
