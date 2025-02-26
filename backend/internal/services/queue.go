// internal/services/queue.go
package services

import (
	"encoding/json"
	"fmt"
	"log"

	"backend/internal/models"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

const (
	// Queue names
	ScanQueueName     = "scan_queue"
	CancelQueueName   = "cancel_queue"
	FindingsQueueName = "findings_queue"
	StatusQueueName   = "status_queue"

	// Exchange name
	ExchangeName = "scanner_exchange"

	// Routing keys
	ScanRoutingKey     = "scan"
	CancelRoutingKey   = "cancel"
	FindingsRoutingKey = "findings"
	StatusRoutingKey   = "status"
)

// ScanRequest represents a scan job to be queued
type ScanRequest struct {
	ScanID      uuid.UUID       `json:"scan_id"`
	ScannerType string          `json:"scanner_type"`
	Targets     []models.Target `json:"targets"`
	Parameters  models.JSONB    `json:"parameters"`
}

// StatusUpdate represents a scan status update
type StatusUpdate struct {
	ScanID  uuid.UUID     `json:"scan_id"`
	Status  models.Status `json:"status"`
	Message string        `json:"message,omitempty"`
}

// QueueService handles interactions with RabbitMQ
type QueueService struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewQueueService creates a new queue service
func NewQueueService(amqpURL string) (*QueueService, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	service := &QueueService{
		connection: conn,
		channel:    ch,
	}

	// Initialize exchange and queues
	err = service.setupExchangesAndQueues()
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to setup exchanges and queues: %w", err)
	}

	return service, nil
}

// Close closes the connection and channel
func (s *QueueService) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.connection != nil {
		s.connection.Close()
	}
}

// setupExchangesAndQueues initializes the RabbitMQ exchange and queues
func (s *QueueService) setupExchangesAndQueues() error {
	// Declare the exchange
	err := s.channel.ExchangeDeclare(
		ExchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare and bind all queues
	queues := []struct {
		name       string
		routingKey string
	}{
		{ScanQueueName, ScanRoutingKey},
		{CancelQueueName, CancelRoutingKey},
		{FindingsQueueName, FindingsRoutingKey},
		{StatusQueueName, StatusRoutingKey},
	}

	for _, q := range queues {
		// Declare queue
		_, err = s.channel.QueueDeclare(
			q.name, // name
			true,   // durable
			false,  // delete when unused
			false,  // exclusive
			false,  // no-wait
			nil,    // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", q.name, err)
		}

		// Bind queue to exchange
		err = s.channel.QueueBind(
			q.name,       // queue name
			q.routingKey, // routing key
			ExchangeName, // exchange
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", q.name, err)
		}
	}

	return nil
}

// QueueScan queues a scan request
func (s *QueueService) QueueScan(scanRequest ScanRequest) error {
	// Convert request to JSON
	body, err := json.Marshal(scanRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal scan request: %w", err)
	}

	// Publish to exchange
	err = s.channel.Publish(
		ExchangeName,   // exchange
		ScanRoutingKey, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish scan request: %w", err)
	}

	log.Printf("Queued scan request: %s", scanRequest.ScanID)
	return nil
}

// CancelScan sends a cancellation request for a scan
func (s *QueueService) CancelScan(scanID uuid.UUID) error {
	// Create message with scan ID
	message := map[string]interface{}{
		"scan_id": scanID.String(),
	}

	// Convert to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal cancel request: %w", err)
	}

	// Publish to exchange
	err = s.channel.Publish(
		ExchangeName,     // exchange
		CancelRoutingKey, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish cancel request: %w", err)
	}

	log.Printf("Sent cancellation for scan: %s", scanID)
	return nil
}

// UpdateScanStatus publishes a status update for a scan
func (s *QueueService) UpdateScanStatus(scanID uuid.UUID, status models.Status, message string) error {
	// Create status update message
	statusUpdate := StatusUpdate{
		ScanID:  scanID,
		Status:  status,
		Message: message,
	}

	// Convert to JSON
	body, err := json.Marshal(statusUpdate)
	if err != nil {
		return fmt.Errorf("failed to marshal status update: %w", err)
	}

	// Publish to exchange
	err = s.channel.Publish(
		ExchangeName,     // exchange
		StatusRoutingKey, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish status update: %w", err)
	}

	log.Printf("Published status update for scan %s: %s", scanID, status)
	return nil
}

// PublishFinding publishes a finding to the findings queue
func (s *QueueService) PublishFinding(finding models.Finding) error {
	// Convert finding to JSON
	body, err := json.Marshal(finding)
	if err != nil {
		return fmt.Errorf("failed to marshal finding: %w", err)
	}

	// Publish to exchange
	err = s.channel.Publish(
		ExchangeName,       // exchange
		FindingsRoutingKey, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish finding: %w", err)
	}

	log.Printf("Published finding for scan: %s, target: %s", finding.ScanID, finding.TargetID)
	return nil
}

// ConsumeScanRequests sets up a consumer for scan requests
func (s *QueueService) ConsumeScanRequests(handler func(ScanRequest) error) error {
	// Consume messages from queue
	msgs, err := s.channel.Consume(
		ScanQueueName, // queue
		"",            // consumer (empty = auto-generated)
		false,         // auto-ack (false = manual ack)
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Process messages in a goroutine
	go func() {
		for d := range msgs {
			// Parse message
			var req ScanRequest
			err := json.Unmarshal(d.Body, &req)
			if err != nil {
				log.Printf("Error unmarshaling scan request: %v", err)
				d.Nack(false, true) // Requeue the message
				continue
			}

			// Handle the request
			err = handler(req)
			if err != nil {
				log.Printf("Error handling scan request: %v", err)
				d.Nack(false, true) // Requeue the message
			} else {
				d.Ack(false) // Acknowledge the message (successfully processed)
			}
		}
	}()

	log.Println("Started consuming scan requests")
	return nil
}

// ConsumeCancellationRequests sets up a consumer for cancellation requests
func (s *QueueService) ConsumeCancellationRequests(handler func(uuid.UUID) error) error {
	// Consume messages from queue
	msgs, err := s.channel.Consume(
		CancelQueueName, // queue
		"",              // consumer (empty = auto-generated)
		false,           // auto-ack (false = manual ack)
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Process messages in a goroutine
	go func() {
		for d := range msgs {
			// Parse message
			var req map[string]interface{}
			err := json.Unmarshal(d.Body, &req)
			if err != nil {
				log.Printf("Error unmarshaling cancel request: %v", err)
				d.Nack(false, true) // Requeue the message
				continue
			}

			// Extract scan ID
			scanIDStr, ok := req["scan_id"].(string)
			if !ok {
				log.Printf("Invalid scan_id in cancel request")
				d.Nack(false, false) // Don't requeue the message (bad format)
				continue
			}

			// Parse UUID
			scanID, err := uuid.Parse(scanIDStr)
			if err != nil {
				log.Printf("Error parsing scan ID: %v", err)
				d.Nack(false, false) // Don't requeue the message (bad UUID)
				continue
			}

			// Handle the request
			err = handler(scanID)
			if err != nil {
				log.Printf("Error handling cancel request: %v", err)
				d.Nack(false, true) // Requeue the message
			} else {
				d.Ack(false) // Acknowledge the message (successfully processed)
			}
		}
	}()

	log.Println("Started consuming cancellation requests")
	return nil
}

// ConsumeStatusUpdates sets up a consumer for status updates
func (s *QueueService) ConsumeStatusUpdates(handler func(StatusUpdate) error) error {
	// Consume messages from queue
	msgs, err := s.channel.Consume(
		StatusQueueName, // queue
		"",              // consumer (empty = auto-generated)
		false,           // auto-ack (false = manual ack)
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Process messages in a goroutine
	go func() {
		for d := range msgs {
			// Parse message
			var update StatusUpdate
			err := json.Unmarshal(d.Body, &update)
			if err != nil {
				log.Printf("Error unmarshaling status update: %v", err)
				d.Nack(false, true) // Requeue the message
				continue
			}

			// Handle the update
			err = handler(update)
			if err != nil {
				log.Printf("Error handling status update: %v", err)
				d.Nack(false, true) // Requeue the message
			} else {
				d.Ack(false) // Acknowledge the message (successfully processed)
			}
		}
	}()

	log.Println("Started consuming status updates")
	return nil
}

// ConsumeFindings sets up a consumer for findings
func (s *QueueService) ConsumeFindings(handler func(models.Finding) error) error {
	// Consume messages from queue
	msgs, err := s.channel.Consume(
		FindingsQueueName, // queue
		"",                // consumer (empty = auto-generated)
		false,             // auto-ack (false = manual ack)
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Process messages in a goroutine
	go func() {
		for d := range msgs {
			// Parse message
			var finding models.Finding
			err := json.Unmarshal(d.Body, &finding)
			if err != nil {
				log.Printf("Error unmarshaling finding: %v", err)
				d.Nack(false, true) // Requeue the message
				continue
			}

			// Handle the finding
			err = handler(finding)
			if err != nil {
				log.Printf("Error handling finding: %v", err)
				d.Nack(false, true) // Requeue the message
			} else {
				d.Ack(false) // Acknowledge the message (successfully processed)
			}
		}
	}()

	log.Println("Started consuming findings")
	return nil
}

