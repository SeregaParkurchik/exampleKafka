package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
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

func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Printf("ERROR: Ошибка при парсинге DSN: %v", err)
		return nil, fmt.Errorf("ошибка при парсинге строки подключения: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Printf("ERROR: Ошибка при создании пула соединений: %v", err)
		return nil, fmt.Errorf("ошибка при создании пула соединений: %w", err)
	}

	return pool, nil
}
