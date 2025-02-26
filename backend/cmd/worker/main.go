// cmd/worker/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/scanner"
	"backend/internal/services"
	"backend/internal/worker"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Parse command line flags
	workerID := flag.String("id", "", "Worker ID (optional, random UUID will be generated if not provided)")
	rabbitMQURL := flag.String("rabbitmq", "", "RabbitMQ URL (falls back to env var RABBITMQ_URL)")
	flag.Parse()

	// Generate worker ID if not provided
	if *workerID == "" {
		*workerID = fmt.Sprintf("worker-%s", uuid.New().String()[:8])
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize RabbitMQ connection
	rabbitURL := *rabbitMQURL
	if rabbitURL == "" {
		rabbitURL = os.Getenv("RABBITMQ_URL")
		if rabbitURL == "" {
			rabbitURL = "amqp://guest:guest@localhost:5672/"
		}
	}

	log.Printf("Starting worker %s, connecting to RabbitMQ at %s", *workerID, rabbitURL)

	// Create queue service
	queueService, err := services.NewQueueService(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialize queue service: %v", err)
	}
	defer queueService.Close()

	// Initialize scanner registry
	scannerRegistry := scanner.NewRegistry()

	// Register available scanners
	pingScanner := scanner.NewPingScanner()
	scannerRegistry.Register("ping", pingScanner)

	// Create and start the worker
	scanWorker := worker.NewWorker(queueService, scannerRegistry, *workerID)
	err = scanWorker.Start()
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	// Wait for termination signal
	log.Printf("Worker %s running, press Ctrl+C to exit", *workerID)

	// Create a channel to listen for OS signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	<-signals

	log.Printf("Worker %s shutting down...", *workerID)

	// Allow time for in-flight operations to complete
	log.Println("Allowing time for in-flight operations to complete...")
	time.Sleep(2 * time.Second)

	log.Printf("Worker %s stopped", *workerID)
}
