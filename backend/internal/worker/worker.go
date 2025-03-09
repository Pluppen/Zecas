// internal/worker/worker.go
package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/internal/models"
	"backend/internal/scanner"
	"backend/internal/services"

	"github.com/google/uuid"
)

// Worker manages scan jobs
type Worker struct {
	queueService    *services.QueueService
	scannerRegistry *scanner.Registry
	targetService   *services.TargetService
	serviceService  *services.ServiceService
	activeScans     map[uuid.UUID]context.CancelFunc
	workerID        string
}

// NewWorker creates a new worker
func NewWorker(
	queueService *services.QueueService,
	scannerRegistry *scanner.Registry,
	targetService *services.TargetService,
	serviceService *services.ServiceService,
	workerID string,
) *Worker {
	return &Worker{
		queueService:    queueService,
		scannerRegistry: scannerRegistry,
		targetService:   targetService,
		serviceService:  serviceService,
		activeScans:     make(map[uuid.UUID]context.CancelFunc),
		workerID:        workerID,
	}
}

// Start begins processing scan requests
func (w *Worker) Start() error {
	// Log worker startup
	log.Printf("Worker %s starting...", w.workerID)

	// Listen for cancellation requests
	err := w.queueService.ConsumeCancellationRequests(w.handleCancellation)
	if err != nil {
		return fmt.Errorf("failed to set up cancellation consumer: %w", err)
	}

	// Setup target consumer
	err = w.queueService.ConsumeTargets(w.handleNewTarget)
	if err != nil {
		return fmt.Errorf("failed to set up target consumer: %w", err)
	}

	// Setup relation consumer
	err = w.queueService.ConsumeTargetRelations(w.handleNewRelation)
	if err != nil {
		return fmt.Errorf("failed to set up relation consumer: %w", err)
	}

	// Setup service consumer
	err = w.queueService.ConsumeServices(w.handleNewService)
	if err != nil {
		return fmt.Errorf("failed to set up service consumer: %w", err)
	}

	// Start consuming scan requests
	err = w.queueService.ConsumeScanRequests(w.handleScanRequest)
	if err != nil {
		return fmt.Errorf("failed to set up scan request consumer: %w", err)
	}

	log.Printf("Worker %s ready and waiting for scan requests", w.workerID)
	return nil
}

// handleScanRequest processes a scan request
func (w *Worker) handleScanRequest(request services.ScanRequest) error {
	log.Printf("[Worker %s] Processing scan request %s (type: %s)",
		w.workerID, request.ScanID, request.ScannerType)

	// Update status to running
	err := w.queueService.UpdateScanStatus(
		request.ScanID,
		models.StatusRunning,
		fmt.Sprintf("Started %s scan on worker %s", request.ScannerType, w.workerID),
	)
	if err != nil {
		log.Printf("[Worker %s] Failed to update scan status: %v", w.workerID, err)
	}

	// Get the scanner
	s, err := w.scannerRegistry.Get(request.ScannerType)
	if err != nil {
		errMsg := fmt.Sprintf("Scanner not found: %s", request.ScannerType)
		log.Printf("[Worker %s] %s", w.workerID, errMsg)
		w.queueService.UpdateScanStatus(request.ScanID, models.StatusFailed, errMsg)
		return err
	}

	// Initialize scanner
	ctx, cancel := context.WithCancel(context.Background())
	w.activeScans[request.ScanID] = cancel

	defer func() {
		cancel()
		delete(w.activeScans, request.ScanID)
	}()

	err = s.Initialize(ctx)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to initialize scanner: %v", err)
		log.Printf("[Worker %s] %s", w.workerID, errMsg)
		w.queueService.UpdateScanStatus(request.ScanID, models.StatusFailed, errMsg)
		return err
	}

	// Process each target
	var totalFindings int
	var totalNewTargets int
	var totalRelations int
	var totalServices int
	startTime := time.Now()

	// First check if we have services specified for scanning
	if s.SupportsServices() && len(request.Services) > 0 {
		for i, service := range request.Services {
			// Skip if context cancelled
			if ctx.Err() != nil {
				log.Printf("[Worker %s] Scan %s was cancelled", w.workerID, request.ScanID)
				return nil
			}

			// Update status
			statusMsg := fmt.Sprintf("Scanning service %d/%d: %s:%d",
				i+1, len(request.Services), service.ServiceName, service.Port)
			w.queueService.UpdateScanStatus(request.ScanID, models.StatusRunning, statusMsg)
			log.Printf("[Worker %s] %s", w.workerID, statusMsg)

			// Convert service to scanner format
			scanTarget := s.ConvertService(service)
			if scanTarget == nil {
				continue // Skip if scanner doesn't support this service
			}

			// Run the scan with timeout
			scanCtx, scanCancel := context.WithTimeout(ctx, 5*time.Minute)
			results, err := s.Scan(scanCtx, scanTarget, request.Parameters)
			scanCancel()

			if err != nil {
				if ctx.Err() == context.Canceled {
					log.Printf("[Worker %s] Scan %s was cancelled", w.workerID, request.ScanID)
					return nil
				}
				log.Printf("[Worker %s] Error scanning service %s:%d: %v",
					w.workerID, service.ServiceName, service.Port, err)
				continue
			}

			// Process scan results
			w.processScanResults(results, request.ScanID, service.TargetID, &service.ID)

			totalFindings += len(results.Findings)
			totalNewTargets += len(results.NewTargets)
			totalRelations += len(results.TargetRelations)
			totalServices += len(results.Services)
		}
	}

	// Then process regular targets
	for i, target := range request.Targets {
		// Skip if context cancelled
		if ctx.Err() != nil {
			log.Printf("[Worker %s] Scan %s was cancelled", w.workerID, request.ScanID)
			return nil
		}

		// Skip if scanner doesn't support this target type
		if !s.SupportsTargetType(target.TargetType) {
			log.Printf("[Worker %s] Scanner %s doesn't support target type %s, skipping",
				w.workerID, request.ScannerType, target.TargetType)
			continue
		}

		// Update status
		statusMsg := fmt.Sprintf("Scanning target %d/%d: %s",
			i+1, len(request.Targets), target.Value)
		w.queueService.UpdateScanStatus(request.ScanID, models.StatusRunning, statusMsg)
		log.Printf("[Worker %s] %s", w.workerID, statusMsg)

		// Convert target to scanner format
		scanTarget := s.ConvertTarget(target)
		if scanTarget == nil {
			continue // Skip if conversion fails
		}

		// Run the scan with timeout
		scanCtx, scanCancel := context.WithTimeout(ctx, 5*time.Minute)
		results, err := s.Scan(scanCtx, scanTarget, request.Parameters)
		scanCancel()

		if err != nil {
			if ctx.Err() == context.Canceled {
				log.Printf("[Worker %s] Scan %s was cancelled", w.workerID, request.ScanID)
				return nil
			}
			log.Printf("[Worker %s] Error scanning target %s: %v", w.workerID, target.Value, err)
			continue
		}

		// Process scan results
		w.processScanResults(results, request.ScanID, target.ID, nil)

		totalFindings += len(results.Findings)
		totalNewTargets += len(results.NewTargets)
		totalRelations += len(results.TargetRelations)
		totalServices += len(results.Services)
	}

	// Update status to completed
	duration := time.Since(startTime).Round(time.Millisecond)
	resultMsg := fmt.Sprintf("Completed %s scan in %s. Found: %d findings, %d new targets, %d relations, %d services",
		request.ScannerType, duration, totalFindings, totalNewTargets, totalRelations, totalServices)
	err = w.queueService.UpdateScanStatus(request.ScanID, models.StatusCompleted, resultMsg)
	if err != nil {
		log.Printf("[Worker %s] Failed to update scan status: %v", w.workerID, err)
	}

	log.Printf("[Worker %s] Scan %s completed successfully in %s",
		w.workerID, request.ScanID, duration)
	return nil
}

// processScanResults handles the results of a scan
func (w *Worker) processScanResults(results *models.ScanResults, scanID uuid.UUID, targetID uuid.UUID, serviceID *uuid.UUID) {
	// Get project ID from target
	target, err := w.targetService.GetByID(targetID)
	if err != nil {
		log.Printf("Error retrieving target %s: %v", targetID, err)
		return
	}
	projectID := target.ProjectID

	// Process findings
	for i := range results.Findings {
		// Set scan and target IDs for findings
		results.Findings[i].ScanID = &scanID
		results.Findings[i].TargetID = targetID

		// Set service ID if applicable
		if serviceID != nil {
			results.Findings[i].ServiceID = serviceID
		}

		// Queue finding
		err := w.queueService.PublishFinding(results.Findings[i])
		if err != nil {
			log.Printf("Error publishing finding: %v", err)
		}
	}

	// Process new targets
	for i := range results.NewTargets {
		// Set project ID for new targets
		results.NewTargets[i].ProjectID = projectID

		// Queue new target
		err := w.queueService.PublishTarget(results.NewTargets[i])
		if err != nil {
			log.Printf("Error publishing new target: %v", err)
		}
	}

	// Process target relations
	for i := range results.TargetRelations {
		// If source ID is not set, use the current target ID
		if results.TargetRelations[i].SourceID == uuid.Nil {
			results.TargetRelations[i].SourceID = targetID
		}

		// Queue target relation
		err := w.queueService.PublishTargetRelation(results.TargetRelations[i])
		if err != nil {
			log.Printf("Error publishing target relation: %v", err)
		}
	}

	// Process services
	for i := range results.Services {
		// If target ID is not set, use the current target ID
		if results.Services[i].TargetID == uuid.Nil {
			results.Services[i].TargetID = targetID
		}

		// Queue service
		err := w.queueService.PublishService(results.Services[i])
		if err != nil {
			log.Printf("Error publishing service: %v", err)
		}
	}
}

// handleCancellation handles a cancellation request
func (w *Worker) handleCancellation(scanID uuid.UUID) error {
	log.Printf("[Worker %s] Received cancellation request for scan: %s", w.workerID, scanID)

	// If we have a running scan with this ID, cancel it
	if cancel, exists := w.activeScans[scanID]; exists {
		cancel()
		delete(w.activeScans, scanID)
		log.Printf("[Worker %s] Cancelled scan: %s", w.workerID, scanID)
	} else {
		log.Printf("[Worker %s] No active scan found with ID: %s", w.workerID, scanID)
	}

	return nil
}

// handleNewTarget processes a new target
func (w *Worker) handleNewTarget(target models.Target) error {
	log.Printf("[Worker %s] Received new target: %s (%s)", w.workerID, target.Value, target.TargetType)

	// Create the target
	err := w.targetService.Create(&target)
	if err != nil {
		return fmt.Errorf("failed to create target: %w", err)
	}

	return nil
}

// handleNewRelation processes a new target relation
func (w *Worker) handleNewRelation(relation models.TargetRelation) error {
	log.Printf("[Worker %s] Received new relation: %s -> %s (%s)",
		w.workerID, relation.SourceID, relation.DestinationID, relation.RelationType)

	// Create the relation
	err := w.targetService.CreateRelation(&relation)
	if err != nil {
		return fmt.Errorf("failed to create target relation: %w", err)
	}

	return nil
}

// handleNewService processes a new service
func (w *Worker) handleNewService(service models.Service) error {
	log.Printf("[Worker %s] Received new service: %s:%d on target %s",
		w.workerID, service.ServiceName, service.Port, service.TargetID)

	// Create the service
	err := w.serviceService.Create(&service)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	return nil
}
