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

// DNS Record Type enum values
const (
	RecordTypeA     = "A"
	RecordTypeAAAA  = "AAAA"
	RecordTypeCNAME = "CNAME"
	RecordTypeANAME = "ANAME"
	RecordTypeSOA   = "SOA"
	RecordTypeNS    = "NS"
	RecordTypeMX    = "MX"
	RecordTypeTXT   = "TXT"
)

// TargetType enum values
const (
	TargetTypeIP     = "ip"
	TargetTypeCIDR   = "cidr"
	TargetTypeDomain = "domain"
)

// TargetRelationType enum values
const (
	RelationResolvesTo   = "resolves_to"
	RelationParentOf     = "parent_of"
	RelationChildOf      = "child_of"
	RelationHostsService = "hosts_service"
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
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ProjectID   uuid.UUID        `json:"project_id" gorm:"type:uuid;not null"`
	TargetType  string           `json:"target_type" gorm:"type:varchar(20);not null;check:target_type IN ('ip', 'cidr', 'domain', 'subdomain')"`
	Value       string           `json:"value" gorm:"type:text;not null"`
	Metadata    JSONB            `json:"metadata" gorm:"type:jsonb;default:'{}'::jsonb"`
	Findings    []Finding        `json:"findings,omitempty" gorm:"foreignKey:TargetID;constraint:OnDelete:CASCADE;"`
	Services    []Service        `json:"services,omitempty" gorm:"foreignKey:TargetID;constraint:OnDelete:CASCADE;"`
	RelatedFrom []TargetRelation `json:"related_from,omitempty" gorm:"foreignKey:SourceID;constraint:OnDelete:CASCADE;"`
	RelatedTo   []TargetRelation `json:"related_to,omitempty" gorm:"foreignKey:DestinationID;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time        `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time        `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TargetRelation represents a relationship between two targets
type TargetRelation struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SourceID      uuid.UUID `json:"source_id" gorm:"type:uuid;not null"`
	DestinationID uuid.UUID `json:"destination_id" gorm:"type:uuid;not null"`
	RelationType  string    `json:"relation_type" gorm:"type:varchar(50);not null"`
	Metadata      JSONB     `json:"metadata" gorm:"type:jsonb;default:'{}'::jsonb"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	Source      Target `json:"-" gorm:"foreignKey:SourceID;constraint:OnDelete:CASCADE"`
	Destination Target `json:"-" gorm:"foreignKey:DestinationID;constraint:OnDelete:CASCADE"`
}

// Service represents a service running on a target
type Service struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TargetID    uuid.UUID `json:"target_id" gorm:"type:uuid;not null"`
	Port        int       `json:"port" gorm:"not null"`
	Protocol    string    `json:"protocol" gorm:"type:varchar(20);not null"`
	ServiceName string    `json:"service_name" gorm:"type:varchar(100)"`
	Version     string    `json:"version" gorm:"type:varchar(100)"`
	Title       string    `json:"title" gorm:"type:varchar(255)"`
	Description string    `json:"description" gorm:"type:text"`
	Banner      string    `json:"banner" gorm:"type:text"`
	RawInfo     JSONB     `json:"raw_info" gorm:"type:jsonb;default:'{}'::jsonb"`
	Findings    []Finding `json:"findings,omitempty" gorm:"foreignKey:ServiceID"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	Target Target `json:"-" gorm:"foreignKey:TargetID"`
}

// Application represents an high level application running o a target
type Application struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ProjectID   uuid.UUID  `json:"project_id" gorm:"type:uuid;not null"`
	ScanID      *uuid.UUID `json:"scan_id,omitempty" gorm:"type:uuid;"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null"`
	Type        string     `json:"type" gorm:"type:varchar(100);not null"` // gitlab, wordpress, jira, etc.
	Version     string     `json:"version" gorm:"type:varchar(100)"`
	Description string     `json:"description" gorm:"type:text"`
	URL         string     `json:"url" gorm:"type:text"`
	HostTarget  *uuid.UUID `json:"host_target,omitempty" gorm:"type:uuid"` // Optional link to host target
	ServiceID   *uuid.UUID `json:"service_id,omitempty" gorm:"type:uuid"`  // Optional link to hosting service
	Metadata    JSONB      `json:"metadata" gorm:"type:jsonb;default:'{}'::jsonb"`
	Findings    []Finding  `json:"findings,omitempty" gorm:"foreignKey:ApplicationID;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// ScanConfig represents a reusable scan configuration
type ScanConfig struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	ScannerType string    `json:"scanner_type" gorm:"type:varchar(50);not null;check:scanner_type IN ('nmap', 'dns', 'subdomain', 'nuclei', 'httpx')"`
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
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ScanID        *uuid.UUID `json:"scan_id,omitempty" gorm:"type:uuid;"`
	TargetID      uuid.UUID  `json:"target_id" gorm:"type:uuid;not null"`
	ServiceID     *uuid.UUID `json:"service_id,omitempty" gorm:"type:uuid;"`
	ApplicationID *uuid.UUID `json:"application_id,omitempty" gorm:"type:uuid;"`
	Title         string     `json:"title" gorm:"type:varchar(255);not null"`
	Description   string     `json:"description" gorm:"type:text"`
	Severity      string     `json:"severity" gorm:"type:varchar(20);not null;check:severity IN ('critical', 'high', 'medium', 'low', 'info', 'unknown')"`
	FindingType   string     `json:"finding_type" gorm:"type:varchar(50);not null"`
	Details       JSONB      `json:"details" gorm:"type:jsonb;default:'{}'::jsonb"`
	DiscoveredAt  time.Time  `json:"discovered_at" gorm:"default:CURRENT_TIMESTAMP"`
	Verified      bool       `json:"verified" gorm:"default:false"`
	Fixed         bool       `json:"fixed" gorm:"default:false"`
	Manual        bool       `json:"manual" gorm:"default:false"`
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
	ServiceIDs   []uuid.UUID `json:"service_ids,omitempty"`
}

// ScanResults represents the output of a scan with possible new targets and relations
type ScanResults struct {
	Findings        []Finding        `json:"findings"`
	NewTargets      []Target         `json:"new_targets,omitempty"`
	TargetRelations []TargetRelation `json:"target_relations,omitempty"`
	Services        []Service        `json:"services,omitempty"`
	Applications    []Application    `json:"applications,omitempty"`
	DNSRecords      []DNSRecord      `json:"dns_records,omitempty"`
	Certificates    []Certificate    `json:"certificates,omitempty"`
}

type DNSRecord struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ProjectID    uuid.UUID  `json:"project_id" gorm:"type:uuid;"`
	ScanID       *uuid.UUID `json:"scan_id,omitempty" gorm:"type:uuid;"`
	TargetID     uuid.UUID  `json:"target_id" gorm:"type:uuid;not null"`
	RecordType   string     `json:"record_type" gorm:"type:text;not null"`
	RecordValue  string     `json:"record_value" gorm:"type:text"`
	Details      JSONB      `json:"details" gorm:"type:jsonb;default:'{}'::jsonb"`
	DiscoveredAt time.Time  `json:"discovered_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type Certificate struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ScanID        *uuid.UUID `json:"scan_id,omitempty" gorm:"type:uuid;"`
	TargetID      uuid.UUID  `json:"target_id" gorm:"type:uuid;"`
	ServiceID     uuid.UUID  `json:"service_id" gorm:"type:uuid;"`
	ApplicationID uuid.UUID  `json:"application_id" gorm:"type:uuid;"`
	ExpiresAt     time.Time  `json:"expires_at" gorm:"default:CURRENT_TIMESTAMP"`
	IssuedAt      time.Time  `json:"issued_at" gorm:"default:CURRENT_TIMESTAMP"`
	Details       JSONB      `json:"details" gorm:"type:jsonb;default:'{}'::jsonb"`
	Issuer        string     `json:"issuer" gorm:"type:text;"`
	Domain        string     `json:"domain" gorm:"type:text;"`
	DiscoveredAt  time.Time  `json:"discovered_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// Auth Related
// The following models needs to adhere to the following schema https://authjs.dev/getting-started/adapters/pg?framework=next-js#schema

type Tabler interface {
	TableName() string
}

type VerificationToken struct {
	Identifier string    `json:"identifier" gorm:"type:text;primary_key"`
	Expires    time.Time `json:"expires" gorm:"default:CURRENT_TIMESTAMP"`
	Token      string    `json:"token" gorm:"type:text;primary_key"`
}

// TableName overrides the table name used by User to `profiles`
func (VerificationToken) TableName() string {
	return "verification_token"
}

type Account struct {
	Id                string    `json:"id" gorm:"type:serial;primary_key"`
	UserId            int       `json:"userId" gorm:"type:int;column:userId"`
	Type              string    `json:"type" gorm:"type:text"`
	Provider          string    `json:"provider" gorm:"type:text"`
	ProviderAccountId string    `json:"providerAccountId" gorm:"type:text;column:providerAccountId"`
	RefreshToken      string    `json:"refresh_token" gorm:"type:text"`
	AccessToken       string    `json:"access_token" gorm:"type:text"`
	ExpiresAt         time.Time `json:"expires_at" gorm:"default:CURRENT_TIMESTAMP"`
	IdToken           string    `json:"id_token" gorm:"type:text"`
	Scope             string    `json:"scope" gorm:"type:text"`
	SessionState      string    `json:"session_state" gorm:"type:text"`
	TokenType         string    `json:"token_type" gorm:"type:text"`
}

type Session struct {
	Id           string    `json:"id" gorm:"type:serial;primary_key"`
	UserId       int       `json:"userId" gorm:"type:int;column:userId"`
	Expires      time.Time `json:"expires" gorm:"default:CURRENT_TIMESTAMP"`
	SessionToken string    `json:"sessionToken" gorm:"type:text;column:sessionToken"`
}

type User struct {
	Id            string    `json:"id" gorm:"type:serial;primary_key"`
	Name          string    `json:"name" gorm:"type:text"`
	Email         string    `json:"email" gorm:"type:text"`
	EmailVerified time.Time `json:"emailVerified" gorm:"default:CURRENT_TIMESTAMP;column:emailVerified"`
	Image         string    `json:"image" gorm:"type:text"`
}
