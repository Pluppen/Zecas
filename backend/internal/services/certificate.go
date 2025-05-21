// internal/services/application.go
package services

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CertificateService struct {
	db *gorm.DB
}

func NewCertificateService(db *gorm.DB) *CertificateService {
	return &CertificateService{db: db}
}

// GetAll returns all applications
func (s *CertificateService) GetAll() ([]models.Certificate, error) {
	var applications []models.Certificate
	result := s.db.Find(&applications)
	return applications, result.Error
}

// GetByID returns a specific application by ID
func (s *CertificateService) GetByID(id uuid.UUID) (*models.Certificate, error) {
	var application models.Certificate
	result := s.db.First(&application, id)
	return &application, result.Error
}

// GetByProjectID returns all applications for a specific project
func (s *CertificateService) GetByProjectID(projectID uuid.UUID) ([]models.Certificate, error) {
	var applications []models.Certificate
	result := s.db.Where("project_id = ?", projectID).Find(&applications)
	return applications, result.Error
}

// GetByType returns all applications of a specific type
func (s *CertificateService) GetByType(appType string) ([]models.Certificate, error) {
	var applications []models.Certificate
	result := s.db.Where("type = ?", appType).Find(&applications)
	return applications, result.Error
}

// GetByTargetID returns all applications hosted on a specific target
func (s *CertificateService) GetByTargetID(targetID uuid.UUID) ([]models.Certificate, error) {
	var applications []models.Certificate
	result := s.db.Where("host_target = ?", targetID).Find(&applications)
	return applications, result.Error
}

// Create creates a new application
func (s *CertificateService) Create(application *models.Certificate) error {
	return s.db.Create(application).Error
}

// Update updates an existing application
func (s *CertificateService) Update(application *models.Certificate) error {
	return s.db.Save(application).Error
}

// Delete deletes an application
func (s *CertificateService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Certificate{}, id).Error
}

// GetCertificates returns all certificates for a specific application
func (s *CertificateService) GetCertificates(applicationID uuid.UUID) ([]models.Certificate, error) {
	var certificates []models.Certificate
	result := s.db.Where("application_id = ?", applicationID).Find(&certificates)
	return certificates, result.Error
}

// BulkCreate creates multiple applications at once
func (s *CertificateService) BulkCreate(applications []models.Certificate) error {
	if len(applications) == 0 {
		return nil
	}
	return s.db.Create(&applications).Error
}

// CountByProject counts the number of applications for a project
func (s *CertificateService) CountByProject(projectID uuid.UUID) (int64, error) {
	var count int64
	result := s.db.Model(&models.Certificate{}).Where("project_id = ?", projectID).Count(&count)
	return count, result.Error
}

// GetFiltered gets DNS Records by filters
func (s *CertificateService) GetFiltered(projectID *uuid.UUID, recordType string) ([]models.Certificate, error) {
	query := s.db.Model(&models.Certificate{})

	if projectID != nil {
		query = query.Joins("JOIN scans ON dns_record.scan_id = scans.id").Where("scans.project_id = ?", projectID)
	}

	if recordType != "" {
		query = query.Where("record_type = ?", recordType)
	}

	var Certificates []models.Certificate
	result := query.Find(&Certificates)
	return Certificates, result.Error
}
