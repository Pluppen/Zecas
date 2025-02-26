// internal/services/project.go
package services

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectService struct {
	db *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}

// GetAll returns all projects
func (s *ProjectService) GetAll() ([]models.Project, error) {
	var projects []models.Project
	result := s.db.Find(&projects)
	return projects, result.Error
}

// GetByID returns a specific project by ID with its targets
func (s *ProjectService) GetByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	result := s.db.Preload("Targets").First(&project, id)
	return &project, result.Error
}

// Create creates a new project
func (s *ProjectService) Create(project *models.Project) error {
	return s.db.Create(project).Error
}

// Update updates an existing project
func (s *ProjectService) Update(project *models.Project) error {
	return s.db.Save(project).Error
}

// Delete deletes a project
func (s *ProjectService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Project{}, id).Error
}

// GetScans returns all scans for a project
func (s *ProjectService) GetScans(projectID uuid.UUID) ([]models.Scan, error) {
	var scans []models.Scan
	result := s.db.Where("project_id = ?", projectID).Find(&scans)
	return scans, result.Error
}

// GetFindings returns all findings for a project
func (s *ProjectService) GetFindings(projectID uuid.UUID) ([]models.Finding, error) {
	var findings []models.Finding
	result := s.db.Joins("JOIN scans ON findings.scan_id = scans.id").
		Where("scans.project_id = ?", projectID).
		Find(&findings)
	return findings, result.Error
}
