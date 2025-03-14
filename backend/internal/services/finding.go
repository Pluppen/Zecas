package services

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FindingService struct {
	db *gorm.DB
}

func NewFindingService(db *gorm.DB) *FindingService {
	return &FindingService{db: db}
}

// GetAll returns all findings
func (s *FindingService) GetAll() ([]models.Finding, error) {
	var findings []models.Finding
	result := s.db.Find(&findings)
	return findings, result.Error
}

// GetByID returns a specific finding by ID
func (s *FindingService) GetByID(id uuid.UUID) (*models.Finding, error) {
	var finding models.Finding
	result := s.db.First(&finding, id)
	return &finding, result.Error
}

// Create creates a new finding
func (s *FindingService) Create(finding *models.Finding) error {
	return s.db.Create(finding).Error
}

// Update updates an existing finding
func (s *FindingService) Update(finding *models.Finding) error {
	return s.db.Save(finding).Error
}

// Delete deletes a finding
func (s *FindingService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Finding{}, id).Error
}

// GetFiltered returns findings filtered by various criteria
func (s *FindingService) GetFiltered(projectID *uuid.UUID, severity, findingType, fixed string) ([]models.Finding, error) {
	query := s.db.Model(&models.Finding{})

	if projectID != nil {
		// Join with scans to filter by project ID
		query = query.Joins("JOIN scans ON findings.scan_id = scans.id").
			Where("scans.project_id = ?", projectID)
	}

	if severity != "" {
		query = query.Where("severity = ?", severity)
	}

	if findingType != "" {
		query = query.Where("finding_type = ?", findingType)
	}

	if fixed == "true" {
		query = query.Where("fixed = ?", true)
	} else if fixed == "false" {
		query = query.Where("fixed = ?", false)
	}

	var findings []models.Finding
	result := query.Find(&findings)
	return findings, result.Error
}

// BulkUpdate updates multiple findings at once
func (s *FindingService) BulkUpdate(ids []uuid.UUID, fixed *bool, verified *bool) error {
	updates := make(map[string]interface{})

	if fixed != nil {
		updates["fixed"] = *fixed
	}

	if verified != nil {
		updates["verified"] = *verified
	}

	if len(updates) == 0 {
		return nil
	}

	return s.db.Model(&models.Finding{}).Where("id IN ?", ids).Updates(updates).Error
}

// CountBySeverity counts findings by severity for a project
func (s *FindingService) CountBySeverity(projectID uuid.UUID) (map[string]int64, error) {
	type Result struct {
		Severity string
		Count    int64
	}

	var results []Result

	err := s.db.Model(&models.Finding{}).
		Select("severity, count(*) as count").
		Joins("JOIN scans ON findings.scan_id = scans.id").
		Where("scans.project_id = ?", projectID).
		Group("severity").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Severity] = r.Count
	}

	return counts, nil
}

// GetLatestFindings gets the most recent findings
func (s *FindingService) GetLatestFindings(limit int) ([]models.Finding, error) {
	var findings []models.Finding
	result := s.db.Order("discovered_at DESC").Limit(limit).Find(&findings)
	return findings, result.Error
}

// GetFindingsByTarget gets all findings for a specific target
func (s *FindingService) GetFindingsByTarget(targetID uuid.UUID) ([]models.Finding, error) {
	var findings []models.Finding
	result := s.db.Where("target_id = ?", targetID).Find(&findings)
	return findings, result.Error
}

// MarkFixed marks a finding as fixed
func (s *FindingService) MarkFixed(id uuid.UUID, fixed bool) error {
	return s.db.Model(&models.Finding{}).Where("id = ?", id).Update("fixed", fixed).Error
}

// MarkVerified marks a finding as verified
func (s *FindingService) MarkVerified(id uuid.UUID, verified bool) error {
	return s.db.Model(&models.Finding{}).Where("id = ?", id).Update("verified", verified).Error
}

// UpsertFinding creates a finding if it doesn't exist or returns existing one
func (s *FindingService) UpsertFinding(finding *models.Finding) (*models.Finding, error) {
	// Try to find an existing finding with the same project_id, finding_type, and value
	var existingFinding models.Finding
	result := s.db.Where(
		"(target_id = ? OR service_id = ? OR application_id = ?) AND finding_type = ? AND severity = ?",
		finding.TargetID, finding.ServiceID, finding.ApplicationID, finding.FindingType, finding.Severity,
	).First(&existingFinding)

	if result.Error == nil {
		// Finding already exists, update metadata if necessary
		if finding.Details != nil {
			// Merge metadata, preferring existing values but adding new ones
			for k, v := range finding.Details {
				if _, exists := existingFinding.Details[k]; !exists {
					if existingFinding.Details == nil {
						existingFinding.Details = models.JSONB{}
					}
					existingFinding.Details[k] = v
				}
			}
			s.db.Model(&existingFinding).Update("metadata", existingFinding.Details)
		}
		return &existingFinding, nil
	}

	// Finding doesn't exist, create it
	err := s.db.Create(finding).Error
	if err != nil {
		return nil, err
	}

	return finding, nil
}
