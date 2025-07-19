package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"activity_tracker/app-consumer/internal/kafka"
	"activity_tracker/app-consumer/internal/storage"
	"activity_tracker/config"
	models "activity_tracker/pkg"
)

func main() {
	log.Println("Запуск App Consumer...")

	pgDSN := config.LoadDSN()
	log.Println(pgDSN)
	kafkaCfg := config.LoadKafkaConfig()

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	dbPool, err := config.New(appCtx, pgDSN)
	if err != nil {
		log.Fatalf("Ошибка при создании пула соединений PostgreSQL: %v", err)
	}
	pgStore := storage.NewPostgresStorage(dbPool)
	defer pgStore.Close()

	consumer := kafka.NewConsumer(kafkaCfg.Brokers, kafkaCfg.Topic, kafkaCfg.GroupID)
	defer func() {
		log.Println("Закрытие Kafka consumer...")
		if err := consumer.Close(); err != nil {
			log.Printf("Ошибка при закрытии Kafka consumer: %v", err)
		}
	}()
	log.Printf("Kafka consumer инициализирован для топика %s с Group ID %s", kafkaCfg.Topic, kafkaCfg.GroupID)

	go func() {
		for {
			select {
			case <-appCtx.Done():
				log.Println("Горутина обработки сообщений останавливается из-за отмены контекста.")
				return
			default:
				msg, err := consumer.ReadMessage(appCtx)
				if err != nil {
					if err.Error() == "kafka: no messages received" || err.Error() == "context canceled" {
					} else {
						log.Printf("Ошибка чтения сообщения из Kafka: %v", err)
					}
					time.Sleep(500 * time.Millisecond)
					continue
				}

				var event models.UserActivityEvent
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					log.Printf("Ошибка десериализации события из Kafka (Топик: %s, Партиция: %d, Офсет: %d): %v. Значение сообщения: %s",
						msg.Topic, msg.Partition, msg.Offset, err, string(msg.Value))
					continue
				}

				log.Printf("Получено и разобрано событие: UserID=%s, EventType=%s, Timestamp=%s, Data=%v",
					event.UserID, event.EventType, event.Timestamp.Format(time.RFC3339), event.Data)

				if err := pgStore.SaveEvent(appCtx, event); err != nil {
					log.Printf("Ошибка сохранения события в PostgreSQL для UserID=%s, EventType=%s: %v", event.UserID, event.EventType, err)
				} else {
					log.Printf("Событие UserID=%s, EventType=%s успешно сохранено в БД.", event.UserID, event.EventType)
				}
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Получен сигнал завершения работы. Запуск корректного завершения...")
	appCancel()
	time.Sleep(2 * time.Second)
	log.Println("App Consumer завершил работу.")
}
