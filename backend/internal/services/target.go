// internal/services/target.go
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
	result := s.db.Preload("Findings").Preload("Services").Where("project_id = ?", projectID).Find(&targets)
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

// ----- Target Relations Methods -----

// CreateRelation creates a new target relation
func (s *TargetService) CreateRelation(relation *models.TargetRelation) error {
	return s.db.Create(relation).Error
}

// GetRelationByID returns a specific relation by ID
func (s *TargetService) GetRelationByID(id uuid.UUID) (*models.TargetRelation, error) {
	var relation models.TargetRelation
	result := s.db.First(&relation, id)
	return &relation, result.Error
}

// GetRelations returns all relations with optional filtering
func (s *TargetService) GetRelations(sourceID, destinationID *uuid.UUID, relationType string) ([]models.TargetRelation, error) {
	var relations []models.TargetRelation
	query := s.db.Model(&models.TargetRelation{})

	if sourceID != nil {
		query = query.Where("source_id = ?", sourceID)
	}

	if destinationID != nil {
		query = query.Where("destination_id = ?", destinationID)
	}

	if relationType != "" {
		query = query.Where("relation_type = ?", relationType)
	}

	result := query.Find(&relations)
	return relations, result.Error
}

// DeleteRelation deletes a relation
func (s *TargetService) DeleteRelation(id uuid.UUID) error {
	return s.db.Delete(&models.TargetRelation{}, id).Error
}

// GetRelatedTargets returns all targets related to a specific target
func (s *TargetService) GetRelatedTargets(targetID uuid.UUID, relationType string) ([]models.Target, error) {
	var targets []models.Target

	// Get targets where the specified target is the source
	query := s.db.Joins("JOIN target_relations ON targets.id = target_relations.destination_id").
		Where("target_relations.source_id = ?", targetID)

	if relationType != "" {
		query = query.Where("target_relations.relation_type = ?", relationType)
	}

	if err := query.Find(&targets).Error; err != nil {
		return nil, err
	}

	// Get targets where the specified target is the destination
	query = s.db.Joins("JOIN target_relations ON targets.id = target_relations.source_id").
		Where("target_relations.destination_id = ?", targetID)

	if relationType != "" {
		query = query.Where("target_relations.relation_type = ?", relationType)
	}

	var additionalTargets []models.Target
	if err := query.Find(&additionalTargets).Error; err != nil {
		return nil, err
	}

	// Combine results
	targets = append(targets, additionalTargets...)

	return targets, nil
}

// BulkCreateRelations creates multiple target relations at once
func (s *TargetService) BulkCreateRelations(relations []models.TargetRelation) error {
	if len(relations) == 0 {
		return nil
	}
	return s.db.Create(&relations).Error
}

// GetTargetsByValue searches for targets by value pattern
func (s *TargetService) GetTargetsByValue(valuePattern string, targetType string) ([]models.Target, error) {
	var targets []models.Target
	query := s.db.Model(&models.Target{}).Where("value LIKE ?", "%"+valuePattern+"%")

	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}

	result := query.Find(&targets)
	return targets, result.Error
}

// GetServices returns all services for a specific target
func (s *TargetService) GetServices(targetID uuid.UUID) ([]models.Service, error) {
	var services []models.Service
	result := s.db.Where("target_id = ?", targetID).Find(&services)
	return services, result.Error
}

// UpsertTarget creates a target if it doesn't exist or returns existing one
func (s *TargetService) UpsertTarget(target *models.Target) (*models.Target, error) {
	// Try to find an existing target with the same project_id, target_type, and value
	var existingTarget models.Target
	result := s.db.Where(
		"project_id = ? AND target_type = ? AND value = ?",
		target.ProjectID, target.TargetType, target.Value,
	).First(&existingTarget)

	if result.Error == nil {
		// Target already exists, update metadata if necessary
		if target.Metadata != nil {
			// Merge metadata, preferring existing values but adding new ones
			for k, v := range target.Metadata {
				if _, exists := existingTarget.Metadata[k]; !exists {
					if existingTarget.Metadata == nil {
						existingTarget.Metadata = models.JSONB{}
					}
					existingTarget.Metadata[k] = v
				}
			}
			s.db.Model(&existingTarget).Update("metadata", existingTarget.Metadata)
		}
		return &existingTarget, nil
	}

	// Target doesn't exist, create it
	err := s.db.Create(target).Error
	if err != nil {
		return nil, err
	}

	return target, nil
}

// FindByTypeAndValue finds a target by type and value in a specific project
func (s *TargetService) FindByTypeAndValue(projectID uuid.UUID, targetType, value string) (*models.Target, error) {
	var target models.Target
	result := s.db.Where(
		"project_id = ? AND target_type = ? AND value = ?",
		projectID, targetType, value,
	).First(&target)
	return &target, result.Error
}
