package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokerAddresses []string, topic, groupID string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          brokerAddresses,
		Topic:            topic,
		GroupID:          groupID,
		MaxWait:          1 * time.Second,
		ReadBatchTimeout: 5 * time.Second,
		CommitInterval:   time.Second,
		StartOffset:      kafka.FirstOffset,
	})
	log.Printf("Kafka consumer initialized for topic %s with group ID %s from brokers %v", topic, groupID, brokerAddresses)
	return &Consumer{reader: r}
}

func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	m, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, fmt.Errorf("failed to read message from Kafka: %w", err)
	}
	return m, nil
}

func (c *Consumer) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	err := c.reader.CommitMessages(ctx, msgs...)
	if err != nil {
		return fmt.Errorf("failed to commit messages: %w", err)
	}
	return nil
}

func (c *Consumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
