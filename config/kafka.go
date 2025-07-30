package config

import (
	"log"
	"os"
	"strings"
)

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

func LoadKafkaConfig() KafkaConfig {

	kafkaBrokersStr := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")

	if kafkaBrokersStr == "" {
		log.Fatalf("Ошибка: Переменная окружения KAFKA_BROKERS не установлена. Пример: 'kafka1:29091,kafka2:29092,kafka3:29093'")
	}
	if kafkaTopic == "" {
		log.Fatalf("Ошибка: Переменная окружения KAFKA_TOPIC не установлена.")
	}

	if kafkaGroupID == "" {
		log.Fatalf("Внимание: Переменная окружения KAFKA_GROUP_ID не установлена. Используйте её для консьюмеров.")

	}

	return KafkaConfig{
		Brokers: strings.Split(kafkaBrokersStr, ","),
		Topic:   kafkaTopic,
		GroupID: kafkaGroupID,
	}
}
