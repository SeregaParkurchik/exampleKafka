// app-producer/internal/handler/activity.go
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	kafkaproducer "activity_tracker/app-producer/internal/kafka"
	models "activity_tracker/pkg"
)

type ActivityHandler struct {
	KafkaProducer *kafkaproducer.Producer
}

func NewActivityHandler(p *kafkaproducer.Producer) *ActivityHandler {
	return &ActivityHandler{KafkaProducer: p}
}

func (h *ActivityHandler) TrackActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var event models.UserActivityEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal event: %v", err), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = h.KafkaProducer.ProduceMessage(ctx, []byte(event.UserID), eventBytes)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		http.Error(w, "Failed to record activity", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully recorded event for UserID: %s, EventType: %s", event.UserID, event.EventType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Activity recorded"})
}
