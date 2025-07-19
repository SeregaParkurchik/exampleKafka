// app-producer/cmd/producer/main.go
package main

import (
	"log"
	"net/http"

	"activity_tracker/app-producer/internal/handler"
	"activity_tracker/config"

	kafkaproducer "activity_tracker/app-producer/internal/kafka"
)

func main() {
	kafkaCfg := config.LoadKafkaConfig()
	producer := kafkaproducer.NewProducer(kafkaCfg.Brokers, kafkaCfg.Topic)
	defer producer.Close()

	activityHandler := handler.NewActivityHandler(producer)

	httpPort := config.LoadHTTPServerPort()

	listenAddr := ":" + httpPort
	http.HandleFunc("/track_activity", activityHandler.TrackActivity)
	log.Printf("Starting HTTP server on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
