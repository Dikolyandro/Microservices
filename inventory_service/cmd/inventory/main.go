package main

import (
	"context"
	"log"
	"os"
	"github.com/recktt77/Microservices-First-/inventory_service/config"
	"github.com/recktt77/Microservices-First-/inventory_service/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("C:\\Users\\админ\\Desktop\\Micriservices\\inventory_service\\.env")
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	ctx := context.Background()
	// TODO: add telemetry here when the topic of logging will be covered
	log.Println("MONGO_DB_URI:", os.Getenv("MONGO_DB_URI"))
	// Parse config
	cfg, err := config.New()
	if err != nil {
		log.Printf("failed to parse config: %v", err)

		return
	}

	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Println("failed to setup application:", err)

		return
	}

	err = application.Run()
	if err != nil {
		log.Println("failed to run application: ", err)

		return
	}

}