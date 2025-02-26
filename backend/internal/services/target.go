package services

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TargetService struct {
	db *gorm.DB
}

func NewTargetService(db *gorm.DB) *TargetService {
	return &TargetService{db: db}
}

// GetAll returns all targets
func (s *TargetService) GetAll() ([]models.Target, error) {
	var targets []models.Target
	result := s.db.Find(&targets)
	return targets, result.Error
}

// GetByID returns a specific target by ID
func (s *TargetService) GetByID(id uuid.UUID) (*models.Target, error) {
	var target models.Target
	result := s.db.First(&target, id)
	return &target, result.Error
}

// GetByProjectID returns all targets for a specific project
func (s *TargetService) GetByProjectID(projectID uuid.UUID) ([]models.Target, error) {
	var targets []models.Target
	result := s.db.Where("project_id = ?", projectID).Find(&targets)
	return targets, result.Error
}

// Create creates a new target
func (s *TargetService) Create(target *models.Target) error {
	return s.db.Create(target).Error
}

// Update updates an existing target
func (s *TargetService) Update(target *models.Target) error {
	return s.db.Save(target).Error
}

// Delete deletes a target
func (s *TargetService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Target{}, id).Error
}

// GetFindings returns all findings for a specific target
func (s *TargetService) GetFindings(targetID uuid.UUID) ([]models.Finding, error) {
	var findings []models.Finding
	result := s.db.Where("target_id = ?", targetID).Find(&findings)
	return findings, result.Error
}

// GetByType returns all targets of a specific type for a project
func (s *TargetService) GetByType(projectID uuid.UUID, targetType string) ([]models.Target, error) {
	var targets []models.Target
	result := s.db.Where("project_id = ? AND target_type = ?", projectID, targetType).Find(&targets)
	return targets, result.Error
}

// BulkCreate creates multiple targets at once
func (s *TargetService) BulkCreate(targets []models.Target) error {
	return s.db.Create(&targets).Error
}

// CountByProject counts the number of targets for a project
func (s *TargetService) CountByProject(projectID uuid.UUID) (int64, error) {
	var count int64
	result := s.db.Model(&models.Target{}).Where("project_id = ?", projectID).Count(&count)
	return count, result.Error
}
