// cmd/api/main.go
package main

import (
	"log"
	"os"

	"backend/internal/api"
	"backend/internal/database"
	"backend/internal/models"
	"backend/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://scanuser:scanpass@localhost:5432/scandb?sslmode=disable"
	}

	db, err := database.Connect(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// RabbitMQ connection
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	// Initialize services
	queueService, err := services.NewQueueService(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer queueService.Close()

	projectService := services.NewProjectService(db)
	targetService := services.NewTargetService(db)
	scanService := services.NewScanService(db)
	findingService := services.NewFindingService(db)
	authService := services.NewAuthService(db)
	serviceService := services.NewServiceService(db)
	relationService := services.NewRelationService(db)
	aplicationService := services.NewApplicationService(db)
	dnsRecordService := services.NewDNSRecordService(db)
	certificateService := services.NewCertificateService(db)

	// Setup findings consumer
	err = queueService.ConsumeFindings(func(finding models.Finding) error {
		_, e := findingService.UpsertFinding(&finding)
		return e
	})
	if err != nil {
		log.Fatalf("Failed to set up findings consumer: %v", err)
	}

	// Setup status updates consumer
	err = queueService.ConsumeStatusUpdates(func(update services.StatusUpdate) error {
		return scanService.UpdateStatus(update.ScanID, update.Status)
	})
	if err != nil {
		log.Fatalf("Failed to set up status updates consumer: %v", err)
	}

	err = queueService.ConsumeTargets(func(target models.Target) error {
		return targetService.Create(&target)
	})
	if err != nil {
		log.Fatalf("Failed to set up targets consumer: %v", err)
	}

	// Setup relations consumer
	err = queueService.ConsumeTargetRelations(func(relation models.TargetRelation) error {
		return relationService.Create(&relation)
	})
	if err != nil {
		log.Fatalf("Failed to set up relations consumer: %v", err)
	}

	// Setup services consumer
	err = queueService.ConsumeServices(func(service models.Service) error {
		return serviceService.Create(&service)
	})
	if err != nil {
		log.Fatalf("Failed to set up services consumer: %v", err)
	}

	// Setup router
	router := api.SetupRouter(projectService, targetService, scanService, findingService, queueService, authService, serviceService, relationService, aplicationService, dnsRecordService, certificateService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create and start the server
	server := api.NewServer(router, port)

	// Start the server (with or without TLS)
	var serverErr error
	log.Println("Starting server without TLS...")
	serverErr = server.Start()

	if serverErr != nil {
		log.Fatalf("Server error: %v", serverErr)
	}
}
