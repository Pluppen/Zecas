package services

import (
	"backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScanService struct {
	db *gorm.DB
}

func NewScanService(db *gorm.DB) *ScanService {
	return &ScanService{db: db}
}

// GetAll returns all scans
func (s *ScanService) GetAll() ([]models.Scan, error) {
	var scans []models.Scan
	result := s.db.Find(&scans)
	return scans, result.Error
}

// GetByID returns a specific scan by ID
func (s *ScanService) GetByID(id uuid.UUID) (*models.Scan, error) {
	var scan models.Scan
	result := s.db.First(&scan, id)
	return &scan, result.Error
}

// Create creates a new scan
func (s *ScanService) Create(scan *models.Scan) error {
	return s.db.Create(scan).Error
}

// Update updates an existing scan
func (s *ScanService) Update(scan *models.Scan) error {
	return s.db.Save(scan).Error
}

// Delete deletes a scan
func (s *ScanService) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Scan{}, id).Error
}

// GetFindings returns all findings for a specific scan
func (s *ScanService) GetFindings(scanID uuid.UUID) ([]models.Finding, error) {
	var findings []models.Finding
	result := s.db.Where("scan_id = ?", scanID).Find(&findings)
	return findings, result.Error
}

// GetScansByProject returns all scans for a specific project
func (s *ScanService) GetScansByProject(projectID uuid.UUID) ([]models.Scan, error) {
	var scans []models.Scan
	result := s.db.Where("project_id = ?", projectID).Order("created_at DESC").Find(&scans)
	return scans, result.Error
}

// GetScansByStatus returns all scans with a specific status
func (s *ScanService) GetScansByStatus(status models.Status) ([]models.Scan, error) {
	var scans []models.Scan
	result := s.db.Where("status = ?", status).Find(&scans)
	return scans, result.Error
}

// UpdateStatus updates the status of a scan
func (s *ScanService) UpdateStatus(scanID uuid.UUID, status models.Status) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == models.StatusRunning {
		now := time.Now()
		updates["started_at"] = now
	} else if status == models.StatusCompleted || status == models.StatusFailed {
		now := time.Now()
		updates["completed_at"] = now
	}

	return s.db.Model(&models.Scan{}).Where("id = ?", scanID).Updates(updates).Error
}

// ----- Scan Configuration Methods -----

// GetAllScanConfigs returns all scan configurations
func (s *ScanService) GetAllScanConfigs() ([]models.ScanConfig, error) {
	var configs []models.ScanConfig
	result := s.db.Find(&configs)
	return configs, result.Error
}

// GetActiveScanConfigs returns all active scan configurations
func (s *ScanService) GetActiveScanConfigs() ([]models.ScanConfig, error) {
	var configs []models.ScanConfig
	result := s.db.Where("active = ?", true).Find(&configs)
	return configs, result.Error
}

// GetScanConfigByID returns a specific scan configuration by ID
func (s *ScanService) GetScanConfigByID(id uuid.UUID) (*models.ScanConfig, error) {
	var config models.ScanConfig
	result := s.db.First(&config, id)
	return &config, result.Error
}

// CreateScanConfig creates a new scan configuration
func (s *ScanService) CreateScanConfig(config *models.ScanConfig) error {
	return s.db.Create(config).Error
}

// UpdateScanConfig updates an existing scan configuration
func (s *ScanService) UpdateScanConfig(config *models.ScanConfig) error {
	return s.db.Save(config).Error
}

// DeleteScanConfig deletes a scan configuration
func (s *ScanService) DeleteScanConfig(id uuid.UUID) error {
	return s.db.Delete(&models.ScanConfig{}, id).Error
}

// GetScanConfigsByType returns all scan configurations of a specific type
func (s *ScanService) GetScanConfigsByType(scannerType string) ([]models.ScanConfig, error) {
	var configs []models.ScanConfig
	result := s.db.Where("scanner_type = ? AND active = ?", scannerType, true).Find(&configs)
	return configs, result.Error
}

// ----- Scan Tasks Methods -----

// CreateScanTask creates a new scan task
func (s *ScanService) CreateScanTask(task *models.ScanTask) error {
	return s.db.Create(task).Error
}

// GetScanTasks returns all tasks for a specific scan
func (s *ScanService) GetScanTasks(scanID uuid.UUID) ([]models.ScanTask, error) {
	var tasks []models.ScanTask
	result := s.db.Where("scan_id = ?", scanID).Find(&tasks)
	return tasks, result.Error
}

// UpdateScanTaskStatus updates the status of a scan task
func (s *ScanService) UpdateScanTaskStatus(taskID uuid.UUID, status models.Status, result models.JSONB) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if result != nil {
		updates["result"] = result
	}

	return s.db.Model(&models.ScanTask{}).Where("id = ?", taskID).Updates(updates).Error
}
