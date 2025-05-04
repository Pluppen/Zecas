// internal/services/application.go
package services

import (
	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DNSRecordService struct {
	db *gorm.DB
}

func NewDNSRecordService(db *gorm.DB) *DNSRecordService {
	return &DNSRecordService{db: db}
}

// GetAll returns all applications
func (s *DNSRecordService) GetAll() ([]models.DNSRecord, error) {
	var applications []models.DNSRecord
	result := s.db.Find(&applications)
	return applications, result.Error
}

// GetByID returns a specific application by ID
func (s *DNSRecordService) GetByID(id uuid.UUID) (*models.DNSRecord, error) {
	var application models.DNSRecord
	result := s.db.First(&application, id)
	return &application, result.Error
}

// GetByProjectID returns all applications for a specific project
func (s *DNSRecordService) GetByProjectID(projectID uuid.UUID) ([]models.DNSRecord, error) {
	var applications []models.DNSRecord
	result := s.db.Where("project_id = ?", projectID).Find(&applications)
	return applications, result.Error
}

// GetByType returns all applications of a specific type
func (s *DNSRecordService) GetByType(appType string) ([]models.DNSRecord, error) {
	var applications []models.DNSRecord
	result := s.db.Where("type = ?", appType).Find(&applications)
	return applications, result.Error
}

// GetByTargetID returns all applications hosted on a specific target
func (s *DNSRecordService) GetByTargetID(targetID uuid.UUID) ([]models.DNSRecord, error) {
	var applications []models.DNSRecord
	result := s.db.Where("host_target = ?", targetID).Find(&applications)
	return applications, result.Error
}

// Create creates a new application
func (s *DNSRecordService) Create(application *models.DNSRecord) error {
	return s.db.Create(application).Error
}

// Update updates an existing application
func (s *DNSRecordService) Update(application *models.DNSRecord) error {
	return s.db.Save(application).Error
}

// Delete deletes an application
func (s *DNSRecordService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.DNSRecord{}, id).Error
}

// GetDNSRecords returns all dnsRecords for a specific application
func (s *DNSRecordService) GetDNSRecords(applicationID uuid.UUID) ([]models.DNSRecord, error) {
	var dnsRecords []models.DNSRecord
	result := s.db.Where("application_id = ?", applicationID).Find(&dnsRecords)
	return dnsRecords, result.Error
}

// BulkCreate creates multiple applications at once
func (s *DNSRecordService) BulkCreate(applications []models.DNSRecord) error {
	if len(applications) == 0 {
		return nil
	}
	return s.db.Create(&applications).Error
}

// CountByProject counts the number of applications for a project
func (s *DNSRecordService) CountByProject(projectID uuid.UUID) (int64, error) {
	var count int64
	result := s.db.Model(&models.DNSRecord{}).Where("project_id = ?", projectID).Count(&count)
	return count, result.Error
}

// GetFiltered gets DNS Records by filters
func (s *DNSRecordService) GetFiltered(projectID *uuid.UUID, recordType string) ([]models.DNSRecord, error) {
	query := s.db.Model(&models.DNSRecord{})

	if projectID != nil {
		query = query.Joins("JOIN scans ON dns_record.scan_id = scans.id").Where("scans.project_id = ?", projectID)
	}

	if recordType != "" {
		query = query.Where("record_type = ?", recordType)
	}

	var DNSRecords []models.DNSRecord
	result := query.Find(&DNSRecords)
	return DNSRecords, result.Error
}
