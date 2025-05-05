package main

import (
	"context"
	"log"
	"github.com/recktt77/Microservices-First-/order_service/config"
	"github.com/recktt77/Microservices-First-/order_service/internal/app"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("C:\\Users\\админ\\Desktop\\Micriservices\\order_service\\.env")
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	ctx := context.Background()
	log.Println("MONGO_DB_URI:", os.Getenv("MONGO_DB_URI"))
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