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
	activeScans     map[uuid.UUID]context.CancelFunc
	workerID        string
}

// NewWorker creates a new worker
func NewWorker(queueService *services.QueueService, scannerRegistry *scanner.Registry, workerID string) *Worker {
	return &Worker{
		queueService:    queueService,
		scannerRegistry: scannerRegistry,
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
	var allFindings []models.Finding
	startTime := time.Now()

	for i, target := range request.Targets {
		// Skip if context cancelled
		if ctx.Err() != nil {
			log.Printf("[Worker %s] Scan %s was cancelled", w.workerID, request.ScanID)
			return nil
		}

		// Update status
		statusMsg := fmt.Sprintf("Scanning target %d/%d: %s",
			i+1, len(request.Targets), target.Value)
		w.queueService.UpdateScanStatus(request.ScanID, models.StatusRunning, statusMsg)
		log.Printf("[Worker %s] %s", w.workerID, statusMsg)

		// Convert target to scanner format
		scanTarget := s.ConvertTarget(target)

		// Run the scan with timeout
		scanCtx, scanCancel := context.WithTimeout(ctx, 2*time.Minute)
		findings, err := s.Scan(scanCtx, scanTarget, request.Parameters)
		scanCancel()

		if err != nil {
			if ctx.Err() == context.Canceled {
				log.Printf("[Worker %s] Scan %s was cancelled", w.workerID, request.ScanID)
				return nil
			}
			log.Printf("[Worker %s] Error scanning target %s: %v", w.workerID, target.Value, err)
			continue
		}

		// Process findings
		for i := range findings {
			// Set scan and target IDs
			findings[i].ScanID = &request.ScanID
			findings[i].TargetID = target.ID

			// Queue finding
			err = w.queueService.PublishFinding(findings[i])
			if err != nil {
				log.Printf("[Worker %s] Error publishing finding: %v", w.workerID, err)
			}
		}

		// Collect all findings for results summary
		allFindings = append(allFindings, findings...)

		log.Printf("[Worker %s] Completed scan of target %s, found %d findings",
			w.workerID, target.Value, len(findings))
	}

	// Update status to completed
	duration := time.Since(startTime).Round(time.Millisecond)
	resultMsg := fmt.Sprintf("Completed %s scan with %d findings in %s",
		request.ScannerType, len(allFindings), duration)
	err = w.queueService.UpdateScanStatus(request.ScanID, models.StatusCompleted, resultMsg)
	if err != nil {
		log.Printf("[Worker %s] Failed to update scan status: %v", w.workerID, err)
	}

	log.Printf("[Worker %s] Scan %s completed successfully in %s",
		w.workerID, request.ScanID, duration)
	return nil
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
