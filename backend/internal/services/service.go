// internal/services/service.go
package services

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceService struct {
	db *gorm.DB
}

func NewServiceService(db *gorm.DB) *ServiceService {
	return &ServiceService{db: db}
}

// GetAll returns all services
func (s *ServiceService) GetAll() ([]models.Service, error) {
	var services []models.Service
	result := s.db.Find(&services)
	return services, result.Error
}

// GetByID returns a specific service by ID
func (s *ServiceService) GetByID(id uuid.UUID) (*models.Service, error) {
	var service models.Service
	result := s.db.First(&service, id)
	return &service, result.Error
}

// GetByTargetID returns all services for a specific target
func (s *ServiceService) GetByTargetID(targetID uuid.UUID) ([]models.Service, error) {
	var services []models.Service
	result := s.db.Where("target_id = ?", targetID).Find(&services)
	return services, result.Error
}

// GetByPort returns all services running on a specific port
func (s *ServiceService) GetByPort(port int) ([]models.Service, error) {
	var services []models.Service
	result := s.db.Where("port = ?", port).Find(&services)
	return services, result.Error
}

// GetByServiceName returns all services of a specific type
func (s *ServiceService) GetByServiceName(serviceName string) ([]models.Service, error) {
	var services []models.Service
	result := s.db.Where("service_name LIKE ?", "%"+serviceName+"%").Find(&services)
	return services, result.Error
}

// Create creates a new service
func (s *ServiceService) Create(service *models.Service) error {
	return s.db.Create(service).Error
}

// BulkCreate creates multiple services at once
func (s *ServiceService) BulkCreate(services []models.Service) error {
	if len(services) == 0 {
		return nil
	}
	return s.db.Create(&services).Error
}

// Update updates an existing service
func (s *ServiceService) Update(service *models.Service) error {
	return s.db.Save(service).Error
}

// Delete deletes a service
func (s *ServiceService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Service{}, id).Error
}

// GetFindings returns all findings for a specific service
func (s *ServiceService) GetFindings(serviceID uuid.UUID) ([]models.Finding, error) {
	var findings []models.Finding
	result := s.db.Where("service_id = ?", serviceID).Find(&findings)
	return findings, result.Error
}

// CountByTarget counts the number of services for a target
func (s *ServiceService) CountByTarget(targetID uuid.UUID) (int64, error) {
	var count int64
	result := s.db.Model(&models.Service{}).Where("target_id = ?", targetID).Count(&count)
	return count, result.Error
}

// GetByProjectID returns all services for a specific project
func (s *ServiceService) GetByProjectID(projectID uuid.UUID) ([]models.Service, error) {
	var services []models.Service
	result := s.db.Joins("JOIN targets ON services.target_id = targets.id").
		Where("targets.project_id = ?", projectID).
		Find(&services)
	return services, result.Error
}
