package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	models "activity_tracker/pkg"

	"github.com/segmentio/kafka-go"
)

type ConsumerService interface {
	Run(ctx context.Context) error
	Stop()
}

type PostgresSaver interface {
	SaveEvent(ctx context.Context, event models.UserActivityEvent) error
}

type KafkaMessageReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type Service struct {
	kafkaReader KafkaMessageReader
	pgSaver     PostgresSaver
	stopChan    chan struct{}
	doneChan    chan struct{}
}

func NewService(kafkaReader KafkaMessageReader, pgSaver PostgresSaver) *Service {
	return &Service{
		kafkaReader: kafkaReader,
		pgSaver:     pgSaver,
		stopChan:    make(chan struct{}),
		doneChan:    make(chan struct{}),
	}
}

func (s *Service) processMessage(ctx context.Context) bool {
	msg, err := s.kafkaReader.ReadMessage(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("Контекст отменен при чтении сообщения из Kafka. Завершение обработки.")
			return false
		}

		if err.Error() == "kafka: no messages received" {
			log.Println("Сообщений в Kafka нет, продолжаю ожидание...")
			return true
		}

		log.Printf("Критическая ошибка чтения сообщения из Kafka: %v. Пауза перед следующей попыткой.", err)
		time.Sleep(500 * time.Millisecond)
		return true
	}

	var event models.UserActivityEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Ошибка десериализации события из Kafka (Топик: %s, Партиция: %d, Офсет: %d): %v. Значение сообщения: %s",
			msg.Topic, msg.Partition, msg.Offset, err, string(msg.Value))
		return true
	}

	log.Printf("Получено и разобрано событие: UserID=%s, EventType=%s",
		event.UserID, event.EventType)

	if err := s.pgSaver.SaveEvent(ctx, event); err != nil {
		log.Printf("Ошибка сохранения события в PostgreSQL для UserID=%s, EventType=%s: %v", event.UserID, event.EventType, err)
		return true
	}

	log.Printf("Событие UserID=%s, EventType=%s успешно сохранено в БД.", event.UserID, event.EventType)
	if err := s.kafkaReader.CommitMessages(ctx, msg); err != nil {
		log.Printf("Ошибка коммита сообщения Kafka (Топик: %s, Партиция: %d, Офсет: %d): %v",
			msg.Topic, msg.Partition, msg.Offset, err)

	} else {
		log.Printf("Сообщение Kafka успешно закоммичено (Топик: %s, Партиция: %d, Офсет: %d)",
			msg.Topic, msg.Partition, msg.Offset)
	}
	return true
}

func (s *Service) Run(ctx context.Context) error {
	log.Println("Запуск горутины обработки сообщений Kafka...")
	go func() {
		defer close(s.doneChan)
		for {
			select {
			case <-ctx.Done():
				log.Println("Горутина обработки сообщений останавливается из-за отмены контекста.")
				return
			case <-s.stopChan:
				log.Println("Горутина обработки сообщений останавливается по запросу.")
				return
			default:
				if !s.processMessage(ctx) {
					return
				}
			}
		}
	}()
	return nil
}

func (s *Service) Stop() {
	log.Println("Отправка сигнала остановки ConsumerService...")
	close(s.stopChan)
	<-s.doneChan
	log.Println("ConsumerService успешно остановлен.")
}
