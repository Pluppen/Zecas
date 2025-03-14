// internal/services/application.go
package services

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationService struct {
	db *gorm.DB
}

func NewApplicationService(db *gorm.DB) *ApplicationService {
	return &ApplicationService{db: db}
}

// GetAll returns all applications
func (s *ApplicationService) GetAll() ([]models.Application, error) {
	var applications []models.Application
	result := s.db.Find(&applications)
	return applications, result.Error
}

// GetByID returns a specific application by ID
func (s *ApplicationService) GetByID(id uuid.UUID) (*models.Application, error) {
	var application models.Application
	result := s.db.First(&application, id)
	return &application, result.Error
}

// GetByProjectID returns all applications for a specific project
func (s *ApplicationService) GetByProjectID(projectID uuid.UUID) ([]models.Application, error) {
	var applications []models.Application
	result := s.db.Where("project_id = ?", projectID).Find(&applications)
	return applications, result.Error
}

// GetByType returns all applications of a specific type
func (s *ApplicationService) GetByType(appType string) ([]models.Application, error) {
	var applications []models.Application
	result := s.db.Where("type = ?", appType).Find(&applications)
	return applications, result.Error
}

// GetByTargetID returns all applications hosted on a specific target
func (s *ApplicationService) GetByTargetID(targetID uuid.UUID) ([]models.Application, error) {
	var applications []models.Application
	result := s.db.Where("host_target = ?", targetID).Find(&applications)
	return applications, result.Error
}

// GetByServiceID returns all applications running on a specific service
func (s *ApplicationService) GetByServiceID(serviceID uuid.UUID) ([]models.Application, error) {
	var applications []models.Application
	result := s.db.Where("service_id = ?", serviceID).Find(&applications)
	return applications, result.Error
}

// Create creates a new application
func (s *ApplicationService) Create(application *models.Application) error {
	return s.db.Create(application).Error
}

// Update updates an existing application
func (s *ApplicationService) Update(application *models.Application) error {
	return s.db.Save(application).Error
}

// Delete deletes an application
func (s *ApplicationService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Application{}, id).Error
}

// GetFindings returns all findings for a specific application
func (s *ApplicationService) GetFindings(applicationID uuid.UUID) ([]models.Finding, error) {
	var findings []models.Finding
	result := s.db.Where("application_id = ?", applicationID).Find(&findings)
	return findings, result.Error
}

// BulkCreate creates multiple applications at once
func (s *ApplicationService) BulkCreate(applications []models.Application) error {
	if len(applications) == 0 {
		return nil
	}
	return s.db.Create(&applications).Error
}

// CountByProject counts the number of applications for a project
func (s *ApplicationService) CountByProject(projectID uuid.UUID) (int64, error) {
	var count int64
	result := s.db.Model(&models.Application{}).Where("project_id = ?", projectID).Count(&count)
	return count, result.Error
}

// Search searches for applications by name or URL
func (s *ApplicationService) Search(query string) ([]models.Application, error) {
	var applications []models.Application
	result := s.db.Where("name LIKE ? OR url LIKE ?", "%"+query+"%", "%"+query+"%").Find(&applications)
	return applications, result.Error
}
