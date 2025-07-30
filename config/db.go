package config

import (
	"fmt"
	"log"
	"os"
)

func LoadDSN() string {

	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DBNAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatalf("Ошибка: Одна или несколько обязательных переменных окружения PostgreSQL (PG_HOST, PG_PORT, PG_USER, PG_PASSWORD, PG_DBNAME) не установлены.")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return dsn
}
