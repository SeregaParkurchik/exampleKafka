// app-producer/internal/kafka/producer.go
package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokerAddresses []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddresses...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				log.Printf("Не удалось записать сообщения: %v", err)
			}

		},
	}

	log.Printf("Kafka producer инициализирован для топика %s по адресам: %v", topic, brokerAddresses) // Комментарий на русском
	return &Producer{writer: w}
}

func (p *Producer) ProduceMessage(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
