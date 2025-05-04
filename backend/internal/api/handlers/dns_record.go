package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DNSRecordHandler struct {
	dnsRecordService *services.DNSRecordService
}

func NewDNSRecordHandler(dnsRecordService *services.DNSRecordService) *DNSRecordHandler {
	return &DNSRecordHandler{
		dnsRecordService: dnsRecordService,
	}
}

// GetDNSRecords returns all dnsRecords with optional filtering
// @Summary Get all dnsRecords
// @Description Get all dnsRecords with optional filtering
// @Tags dnsRecords
// @Accept json
// @Produce json
// @Param severity query string false "Filter by severity"
// @Param type query string false "Filter by dnsRecord type"
// @Param fixed query string false "Filter by fixed status"
// @Param project_id query string false "Filter by project ID"
// @Success 200 {array} models.DNSRecord
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/dnsRecords [get]
func (h *DNSRecordHandler) GetDNSRecords(c *gin.Context) {
	// Support filtering
	dnsRecordType := c.Query("type")

	// Get project ID from query param if exists
	projectIDStr := c.Query("project_id")
	var projectID *uuid.UUID

	if projectIDStr != "" {
		id, err := uuid.Parse(projectIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
			return
		}
		projectID = &id
	}

	dnsRecords, err := h.dnsRecordService.GetFiltered(projectID, dnsRecordType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dnsRecords"})
		return
	}

	c.JSON(http.StatusOK, dnsRecords)
}

// GetDNSRecord returns a specific dnsRecord by ID
// @Summary Get a dnsRecord
// @Description Get a specific dnsRecord by ID
// @Tags dnsRecords
// @Accept json
// @Produce json
// @Param id path string true "DNSRecord ID"
// @Success 200 {object} models.DNSRecord
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/dnsRecords/{id} [get]
func (h *DNSRecordHandler) GetDNSRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dnsRecord ID format"})
		return
	}

	dnsRecord, err := h.dnsRecordService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DNSRecord not found"})
		return
	}

	c.JSON(http.StatusOK, dnsRecord)
}

// CreateDNSRecord creates a new dnsRecord
// @Summary Create a dnsRecord
// @Description Create a new dnsRecord
// @Tags dnsRecords
// @Accept json
// @Produce json
// @Param dnsRecord body object true "DNSRecord Details"
// @Success 201 {object} models.DNSRecord
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/dnsRecords [post]
func (h *DNSRecordHandler) CreateDNSRecord(c *gin.Context) {
	var input struct {
		ScanID      *string      `json:"scan_id"`
		ProjectID   string       `json:"project_id" binding:"required"`
		TargetID    string       `json:"target_id" binding:"required"`
		RecordType  string       `json:"record_type" binding:"required"`
		RecordValue string       `json:"record_value"`
		Details     models.JSONB `json:"details"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate severity
	if input.RecordType != models.RecordTypeA &&
		input.RecordType != models.RecordTypeAAAA &&
		input.RecordType != models.RecordTypeCNAME &&
		input.RecordType != models.RecordTypeANAME &&
		input.RecordType != models.RecordTypeSOA &&
		input.RecordType != models.RecordTypeNS &&
		input.RecordType != models.RecordTypeTXT &&
		input.RecordType != models.RecordTypeMX {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record type"})
		return
	}

	var scanID *uuid.UUID
	if input.ScanID != nil && *input.ScanID != "" {
		id, err := uuid.Parse(*input.ScanID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan ID format"})
			return
		}
		scanID = &id
	}

	projectID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	targetID, err := uuid.Parse(input.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	dnsRecord := &models.DNSRecord{
		ScanID:      scanID,
		ProjectID:   projectID,
		TargetID:    targetID,
		RecordType:  input.RecordType,
		RecordValue: input.RecordValue,
		Details:     input.Details,
	}

	err = h.dnsRecordService.Create(dnsRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dnsRecord"})
		return
	}

	c.JSON(http.StatusCreated, dnsRecord)
}

// UpdateDNSRecord updates an existing dnsRecord
// @Summary Update a dnsRecord
// @Description Update an existing dnsRecord
// @Tags dnsRecords
// @Accept json
// @Produce json
// @Param id path string true "DNSRecord ID"
// @Param dnsRecord body object true "Updated DNSRecord Details"
// @Success 200 {object} models.DNSRecord
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/dnsRecords/{id} [put]
func (h *DNSRecordHandler) UpdateDNSRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dnsRecord ID format"})
		return
	}

	var input struct {
		ScanID      *string      `json:"scan_id"`
		ProjectID   string       `json:"project_id" binding:"required"`
		TargetID    string       `json:"target_id"`
		RecordType  string       `json:"record_type"`
		RecordValue *string      `json:"record_value"`
		Details     models.JSONB `json:"details"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dnsRecord, err := h.dnsRecordService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DNSRecord not found"})
		return
	}

	// Update fields if provided
	if input.RecordType != models.RecordTypeA &&
		input.RecordType != models.RecordTypeAAAA &&
		input.RecordType != models.RecordTypeCNAME &&
		input.RecordType != models.RecordTypeANAME &&
		input.RecordType != models.RecordTypeSOA &&
		input.RecordType != models.RecordTypeNS &&
		input.RecordType != models.RecordTypeTXT &&
		input.RecordType != models.RecordTypeMX {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record type"})
		return
	}

	if input.RecordValue != nil {
		dnsRecord.RecordValue = *input.RecordValue
	}

	if input.Details != nil {
		dnsRecord.Details = input.Details
	}

	err = h.dnsRecordService.Update(dnsRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update dnsRecord"})
		return
	}

	c.JSON(http.StatusOK, dnsRecord)
}

// DeleteDNSRecord deletes a dnsRecord
// @Summary Delete a dnsRecord
// @Description Delete a dnsRecord
// @Tags dnsRecords
// @Accept json
// @Produce json
// @Param id path string true "DNSRecord ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/dnsRecords/{id} [delete]
func (h *DNSRecordHandler) DeleteDNSRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dnsRecord ID format"})
		return
	}

	err = h.dnsRecordService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete dnsRecord"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DNSRecord deleted successfully"})
}
