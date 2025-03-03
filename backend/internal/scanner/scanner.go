// internal/scanner/scanner.go
package scanner

import (
	"context"
	"fmt"

	"backend/internal/models"

	"github.com/google/uuid"
)

// Scanner interface defines the methods that each scanner must implement
type Scanner interface {
	// Initialize sets up the scanner with any required configuration
	Initialize(ctx context.Context) error

	// ConvertTarget converts a database target model to a scanner-specific target
	ConvertTarget(target models.Target) interface{}

	// Scan performs the actual scanning operation
	Scan(ctx context.Context, target interface{}, params models.JSONB) ([]models.Finding, error)

	// Type returns the scanner type identifier
	Type() string
}

// Registry stores and provides access to scanner implementations
type Registry struct {
	scanners map[string]Scanner
}

// NewRegistry creates a new scanner registry
func NewRegistry() *Registry {
	return &Registry{
		scanners: make(map[string]Scanner),
	}
}

// Register adds a scanner to the registry
func (r *Registry) Register(name string, scanner Scanner) {
	r.scanners[name] = scanner
}

// Get retrieves a scanner by name
func (r *Registry) Get(name string) (Scanner, error) {
	scanner, exists := r.scanners[name]
	if !exists {
		return nil, fmt.Errorf("scanner %s not found", name)
	}
	return scanner, nil
}

// CreateFinding is a helper function to create a finding
func CreateFinding(
	scanID uuid.UUID,
	targetID uuid.UUID,
	title string,
	description string,
	severity string,
	findingType string,
	details models.JSONB,
) models.Finding {
	return models.Finding{
		ScanID:      &scanID,
		TargetID:    targetID,
		Title:       title,
		Description: description,
		Severity:    severity,
		FindingType: findingType,
		Details:     details,
		Verified:    false,
		Fixed:       false,
	}
}
