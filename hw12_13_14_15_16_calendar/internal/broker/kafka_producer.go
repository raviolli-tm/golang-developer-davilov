package broker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
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

// GetProducer creates the Kafka writer. If the broker is not available yet,
// connection attempts are repeated with a linear backoff until max_retries is
// exhausted or the context is cancelled.
func (k *ProducerKafka) GetProducer(ctx context.Context) (*kafka.Writer, error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.producer != nil {
		return k.producer, nil
	}

	maxRetries := k.retryCfg.MaxRetries
	if maxRetries < 1 {
		maxRetries = 1
	}

	baseTime := time.Duration(k.retryCfg.BaseTime) * time.Second
	if baseTime <= 0 {
		baseTime = time.Second
	}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		conn, err := kafka.DialContext(ctx, "tcp", k.producerCfg.BrokersUrl)
		if err == nil {
			_ = conn.Close()

			k.producer = &kafka.Writer{
				Addr:         kafka.TCP(k.producerCfg.BrokersUrl),
				Topic:        k.producerCfg.Topic,
				Balancer:     &kafka.LeastBytes{},
				MaxAttempts:  3,
				BatchTimeout: 10 * time.Millisecond,
			}

			return k.producer, nil
		}

		lastErr = fmt.Errorf("connect to broker: %w", err)

		if attempt == maxRetries {
			break
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(baseTime):
		}
	}

	return nil, fmt.Errorf("broker is not available after %d attempts: %w", maxRetries, lastErr)
}

func (k *ProducerKafka) SendMessage(ctx context.Context, key, value []byte) error {
	writer, err := k.GetProducer(ctx)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	if err := writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (k *ProducerKafka) Close() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.producer == nil {
		return nil
	}

	err := k.producer.Close()
	k.producer = nil

	return err
}
