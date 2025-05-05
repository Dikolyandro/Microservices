package main

import (
	"log"
	grpcServer "statistics_service/internal/adapter/grpc"
	natsHandler "statistics_service/internal/adapter/nats"
	"statistics_service/internal/adapter/storage"
	"statistics_service/internal/usecase"

	"github.com/nats-io/nats.go"
)

func main() {
	// Подключение к MongoDB
	mongoRepo, err := storage.NewMongoRepo("mongodb://localhost:27017", "statisticsdb")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Подключение к NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Drain()

	// Инициализация бизнес-логики
	statsUC := usecase.New(mongoRepo)

	// Подключение к NATS-событиям
	nats := natsHandler.NewNATSHandler(statsUC)
	nats.Subscribe(nc)

	// Запуск gRPC-сервера (асинхронно)
	go grpcServer.NewServer(statsUC)

	log.Println("Statistics Service is running...")

	select {} // Блокировка main
}
