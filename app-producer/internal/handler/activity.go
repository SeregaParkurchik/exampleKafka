package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	kafkaproducer "activity_tracker/app-producer/internal/kafka"
	models "activity_tracker/pkg"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	KafkaProducer *kafkaproducer.Producer
}

func NewActivityHandler(p *kafkaproducer.Producer) *ActivityHandler {
	return &ActivityHandler{KafkaProducer: p}
}

func (h *ActivityHandler) TrackActivity(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Разрешен только метод POST"})
		return
	}

	var event models.UserActivityEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Printf("Ошибка десериализации JSON: %v, запрос: %s", err, c.Request.URL.Path)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON", "details": err.Error()})
		return
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Ошибка сериализации события в JSON для Kafka: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки события"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.KafkaProducer.ProduceMessage(ctx, []byte(event.UserID), eventBytes)
	if err != nil {
		log.Printf("Не удалось записать сообщение в Kafka для UserID: %s, EventType: %s, ошибка: %v",
			event.UserID, event.EventType, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось записать активность"})
		return
	}

	log.Printf("Событие успешно записано в Kafka для UserID: %s, EventType: %s", event.UserID, event.EventType)
	c.JSON(http.StatusOK, gin.H{"status": "успех", "message": "Активность записана"})
}
