package config

import (
	"log"
	"os"
)

func LoadHTTPServerPort() string {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		log.Fatalf("Внимание: Переменная окружения HTTP_PORT не установлена. Используем дефолтный порт: %s", httpPort)
	}
	return httpPort
}
