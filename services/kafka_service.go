// services/kafka_service.go
package services

import (
	"context"
	"time"

	"crypto-exchange/config"

	"github.com/segmentio/kafka-go"
)

// KafkaService encapsulates the Kafka writer.
type KafkaService struct {
	Writer *kafka.Writer
	Topic  string
}

// NewKafkaService initializes the KafkaService.
func NewKafkaService(cfg config.KafkaConfig) *KafkaService {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topics
		Balancer: &kafka.LeastBytes{},
	})

	return &KafkaService{
		Writer: writer,
		Topic:  cfg.Topic,
	}
}

// Publish sends a message to the Kafka topic.
func (k *KafkaService) Publish(message string) error {
	return k.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Time:  time.Now(),
			Key:   []byte("transaction"),
			Value: []byte(message),
		},
	)
}

// Close terminates the Kafka writer.
func (k *KafkaService) Close() error {
	return k.Writer.Close()
}