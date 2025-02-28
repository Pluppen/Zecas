// internal/models/models.go
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// JSONB type for PostgreSQL jsonb columns
type JSONB map[string]interface{}

// Value for implementing driver.Valuer
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan for implementing sql.Scanner
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

// Status type for enum values
type Status string

// ScanStatus enum values
const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// Severity enum values
const (
	SeverityCritical = "critical"
	SeverityHigh     = "high"
	SeverityMedium   = "medium"
	SeverityLow      = "low"
	SeverityInfo     = "info"
	SeverityUnknown  = "unknown"
)

// TargetType enum values
const (
	TargetTypeIP     = "ip"
	TargetTypeCIDR   = "cidr"
	TargetTypeDomain = "domain"
)

// Project represents a scanning project
type Project struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Targets     []Target  `json:"targets,omitempty" gorm:"foreignKey:ProjectID"`
	Scans       []Scan    `json:"scans,omitempty" gorm:"foreignKey:ProjectID"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// Target represents an individual target for scanning (IP, CIDR, or domain)
type Target struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ProjectID  uuid.UUID `json:"project_id" gorm:"type:uuid;not null"`
	TargetType string    `json:"target_type" gorm:"type:varchar(20);not null;check:target_type IN ('ip', 'cidr', 'domain')"`
	Value      string    `json:"value" gorm:"type:text;not null"`
	Metadata   JSONB     `json:"metadata" gorm:"type:jsonb;default:'{}'::jsonb"`
	Findings   []Finding `json:"findings,omitempty" gorm:"foreignKey:TargetID"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// ScanConfig represents a reusable scan configuration
type ScanConfig struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	ScannerType string    `json:"scanner_type" gorm:"type:varchar(50);not null"`
	Parameters  JSONB     `json:"parameters" gorm:"type:jsonb;default:'{}'::jsonb"`
	Active      bool      `json:"active" gorm:"default:true"`
	Scans       []Scan    `json:"scans,omitempty" gorm:"foreignKey:ScanConfigID"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// Scan represents a scan job
type Scan struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ProjectID    uuid.UUID  `json:"project_id" gorm:"type:uuid;not null"`
	ScanConfigID uuid.UUID  `json:"scan_config_id" gorm:"type:uuid;not null"`
	Status       Status     `json:"status" gorm:"type:varchar(50);not null;default:'pending'"`
	StartedAt    *time.Time `json:"started_at" gorm:"type:timestamp with time zone"`
	CompletedAt  *time.Time `json:"completed_at" gorm:"type:timestamp with time zone"`
	RawResults   JSONB      `json:"raw_results" gorm:"type:jsonb"`
	Error        string     `json:"error" gorm:"type:text"`
	Findings     []Finding  `json:"findings,omitempty" gorm:"foreignKey:ScanID"`
	ScanTasks    []ScanTask `json:"scan_tasks,omitempty" gorm:"foreignKey:ScanID"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// Finding represents a vulnerability or other issue found during scanning
type Finding struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ScanID       *uuid.UUID `json:"scan_id,omitempty" gorm:"type:uuid;"`
	TargetID     uuid.UUID  `json:"target_id" gorm:"type:uuid;not null"`
	Title        string     `json:"title" gorm:"type:varchar(255);not null"`
	Description  string     `json:"description" gorm:"type:text"`
	Severity     string     `json:"severity" gorm:"type:varchar(20);not null;check:severity IN ('critical', 'high', 'medium', 'low', 'info', 'unknown')"`
	FindingType  string     `json:"finding_type" gorm:"type:varchar(50);not null"`
	Details      JSONB      `json:"details" gorm:"type:jsonb;default:'{}'::jsonb"`
	DiscoveredAt time.Time  `json:"discovered_at" gorm:"default:CURRENT_TIMESTAMP"`
	Verified     bool       `json:"verified" gorm:"default:false"`
	Fixed        bool       `json:"fixed" gorm:"default:false"`
	Manual       bool       `json:"manual" gorm:"default:false"`
}

// Report represents a generated report for a project
type Report struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ProjectID uuid.UUID `json:"project_id" gorm:"type:uuid;not null"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Format    string    `json:"format" gorm:"type:varchar(20);not null;check:format IN ('pdf', 'html', 'docx', 'json')"`
	Status    string    `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	URL       string    `json:"url" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// ScanTask represents an individual task within a scan
type ScanTask struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ScanID     uuid.UUID `json:"scan_id" gorm:"type:uuid;not null"`
	TaskType   string    `json:"task_type" gorm:"type:varchar(50);not null"`
	Parameters JSONB     `json:"parameters" gorm:"type:jsonb;default:'{}'::jsonb"`
	Status     Status    `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	Result     JSONB     `json:"result" gorm:"type:jsonb"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// CreateProjectInput represents the input for creating a new project
type CreateProjectInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	IPRanges    []string `json:"ip_ranges"`
	CIDRRanges  []string `json:"cidr_ranges"`
	Domains     []string `json:"domains"`
}

// StartScanInput represents the input for starting a new scan
type StartScanInput struct {
	ProjectID    uuid.UUID   `json:"project_id" binding:"required"`
	ScanConfigID uuid.UUID   `json:"scan_config_id" binding:"required"`
	TargetIDs    []uuid.UUID `json:"target_ids,omitempty"`
}
