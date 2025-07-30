package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"activity_tracker/app-consumer/internal/kafka"
	"activity_tracker/app-consumer/internal/service"
	"activity_tracker/app-consumer/internal/storage"
	"activity_tracker/config"
)

func main() {
	log.Println("Запуск App Consumer...")

	pgDSN := config.LoadDSN()
	kafkaCfg := config.LoadKafkaConfig()

	appCtx, appCancel := context.WithCancel(context.Background())

	defer appCancel()

	dbPool, err := storage.New(appCtx, pgDSN)
	if err != nil {
		log.Fatalf("Ошибка при создании пула соединений PostgreSQL: %v", err)
	}
	pgStore := storage.NewPostgresStorage(dbPool)
	defer pgStore.Close()

	kafkaConsumer := kafka.NewConsumer(kafkaCfg.Brokers, kafkaCfg.Topic, kafkaCfg.GroupID)
	defer func() {
		log.Println("Закрытие Kafka consumer...")
		if err := kafkaConsumer.Close(); err != nil {
			log.Printf("Ошибка при закрытии Kafka consumer: %v", err)
		}
	}()
	log.Printf("Kafka consumer инициализирован для топика %s с Group ID %s", kafkaCfg.Topic, kafkaCfg.GroupID)

	consumerService := service.NewService(kafkaConsumer, pgStore)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := consumerService.Run(appCtx); err != nil {
			log.Printf("Ошибка при запуске ConsumerService: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Получен сигнал завершения работы. Запуск корректного завершения...")

	appCancel()
	consumerService.Stop()

	log.Println("App Consumer завершил работу.")
}
