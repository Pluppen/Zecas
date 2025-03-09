package services

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RelationService struct {
	db *gorm.DB
}

func NewRelationService(db *gorm.DB) *RelationService {
	return &RelationService{db: db}
}

// GetAll returns all target relations
func (s *RelationService) GetAll() ([]models.TargetRelation, error) {
	var relations []models.TargetRelation
	result := s.db.Find(&relations)
	return relations, result.Error
}

// GetByID returns a specific target relation by ID
func (s *RelationService) GetByID(id uuid.UUID) (*models.TargetRelation, error) {
	var relation models.TargetRelation
	result := s.db.First(&relation, id)
	return &relation, result.Error
}

// GetFiltered returns relations with optional filtering by source, destination, and type
func (s *RelationService) GetFiltered(sourceID, destinationID *uuid.UUID, relationType string) ([]models.TargetRelation, error) {
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

	var relations []models.TargetRelation
	result := query.Find(&relations)
	return relations, result.Error
}

// Create creates a new target relation
func (s *RelationService) Create(relation *models.TargetRelation) error {
	return s.db.Create(relation).Error
}

// Update updates an existing target relation
func (s *RelationService) Update(relation *models.TargetRelation) error {
	return s.db.Save(relation).Error
}

// Delete deletes a target relation
func (s *RelationService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.TargetRelation{}, id).Error
}

// GetRelationsForTarget gets all relations where a specific target is either source or destination
func (s *RelationService) GetRelationsForTarget(targetID uuid.UUID) ([]models.TargetRelation, error) {
	var relations []models.TargetRelation
	result := s.db.Where("source_id = ? OR destination_id = ?", targetID, targetID).Find(&relations)
	return relations, result.Error
}

// GetSourceTargets gets all targets that are sources for a specific destination
func (s *RelationService) GetSourceTargets(destinationID uuid.UUID, relationType string) ([]models.Target, error) {
	var targets []models.Target
	query := s.db.Joins("JOIN target_relations ON targets.id = target_relations.source_id").
		Where("target_relations.destination_id = ?", destinationID)

	if relationType != "" {
		query = query.Where("target_relations.relation_type = ?", relationType)
	}

	result := query.Find(&targets)
	return targets, result.Error
}

// GetDestinationTargets gets all targets that are destinations for a specific source
func (s *RelationService) GetDestinationTargets(sourceID uuid.UUID, relationType string) ([]models.Target, error) {
	var targets []models.Target
	query := s.db.Joins("JOIN target_relations ON targets.id = target_relations.destination_id").
		Where("target_relations.source_id = ?", sourceID)

	if relationType != "" {
		query = query.Where("target_relations.relation_type = ?", relationType)
	}

	result := query.Find(&targets)
	return targets, result.Error
}

// BulkCreate creates multiple target relations at once
func (s *RelationService) BulkCreate(relations []models.TargetRelation) error {
	if len(relations) == 0 {
		return nil
	}
	return s.db.Create(&relations).Error
}
