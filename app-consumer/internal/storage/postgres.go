package storage

import (
	models "activity_tracker/pkg"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConsumerStorage interface {
	SaveEvent(ctx context.Context, event models.UserActivityEvent) error
	Close()
}
type postgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(pool *pgxpool.Pool) ConsumerStorage {
	return &postgresStorage{pool: pool}
}

func (s *postgresStorage) SaveEvent(ctx context.Context, event models.UserActivityEvent) error {
	jsonData, err := event.DataToJson()
	if err != nil {
		return fmt.Errorf("не удалось сериализовать данные события в JSONB: %w", err)
	}

	query := `INSERT INTO user_activity_events (user_id, event_type, timestamp, data)
              VALUES ($1, $2, $3, $4)`

	_, err = s.pool.Exec(ctx, query, event.UserID, event.EventType, event.Timestamp, jsonData)
	if err != nil {
		return fmt.Errorf("не удалось вставить событие в базу данных: %w", err)
	}
	return nil
}

func (s *postgresStorage) Close() {
	if s.pool != nil {
		log.Println("Закрытие пула соединений PostgreSQL.")
		s.pool.Close()
	}
}
