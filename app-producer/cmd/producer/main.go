package main

import (
	"context"
	"log"
	"net/http" // Нужен для http.ErrServerClosed
	"os"
	"os/signal"
	"syscall"
	"time"

	"activity_tracker/app-producer/internal/handler"
	"activity_tracker/app-producer/internal/kafka"
	"activity_tracker/app-producer/internal/server"
	"activity_tracker/config"
)

func main() {
	log.Println("Запуск App Producer...")

	kafkaCfg := config.LoadKafkaConfig()
	httpPort := config.LoadHTTPServerPort()
	listenAddr := ":" + httpPort

	producer := kafka.NewProducer(kafkaCfg.Brokers, kafkaCfg.Topic)
	defer func() {
		log.Println("Закрытие Kafka producer...")
		if err := producer.Close(); err != nil {
			log.Printf("Ошибка при закрытии Kafka producer: %v", err)
		}
	}()
	log.Printf("Kafka producer инициализирован для топика %s на брокерах %v", kafkaCfg.Topic, kafkaCfg.Brokers)

	activityHandler := handler.NewActivityHandler(producer)

	httpServer := server.NewServer(listenAddr, activityHandler)
	log.Printf("HTTP сервер инициализирован.")

	go func() {
		if err := httpServer.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Получен сигнал завершения работы. Запуск корректного завершения HTTP сервера...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при корректном завершении работы HTTP сервера: %v", err)
	} else {
		log.Println("HTTP сервер успешно остановлен.")
	}

	log.Println("App Producer завершил работу.")
}
