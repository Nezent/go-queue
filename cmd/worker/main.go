package main

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"

	"github.com/Nezent/go-queue/internal/worker"
	"github.com/Nezent/go-queue/internal/worker/processor"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: os.Getenv("REDIS_ADDR"),
	}

	config := processor.SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}

	taskProcessor := processor.NewTaskProcessor(config)

	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
	})

	mux := worker.NewServeMux(taskProcessor)

	log.Println("Starting Asynq worker...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run worker: %v", err)
	}
}
