package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FindingHandler struct {
	findingService *services.FindingService
}

func NewFindingHandler(findingService *services.FindingService) *FindingHandler {
	return &FindingHandler{
		findingService: findingService,
	}
}

// GetFindings returns all findings with optional filtering
// @Summary Get all findings
// @Description Get all findings with optional filtering
// @Tags findings
// @Accept json
// @Produce json
// @Param severity query string false "Filter by severity"
// @Param type query string false "Filter by finding type"
// @Param fixed query string false "Filter by fixed status"
// @Param project_id query string false "Filter by project ID"
// @Success 200 {array} models.Finding
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/findings [get]
func (h *FindingHandler) GetFindings(c *gin.Context) {
	// Support filtering
	severity := c.Query("severity")
	findingType := c.Query("type")
	fixed := c.Query("fixed")

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

	findings, err := h.findingService.GetFiltered(projectID, severity, findingType, fixed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve findings"})
		return
	}

	c.JSON(http.StatusOK, findings)
}

// GetFinding returns a specific finding by ID
// @Summary Get a finding
// @Description Get a specific finding by ID
// @Tags findings
// @Accept json
// @Produce json
// @Param id path string true "Finding ID"
// @Success 200 {object} models.Finding
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/findings/{id} [get]
func (h *FindingHandler) GetFinding(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid finding ID format"})
		return
	}

	finding, err := h.findingService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Finding not found"})
		return
	}

	c.JSON(http.StatusOK, finding)
}

// CreateFinding creates a new finding
// @Summary Create a finding
// @Description Create a new finding
// @Tags findings
// @Accept json
// @Produce json
// @Param finding body object true "Finding Details"
// @Success 201 {object} models.Finding
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/findings [post]
func (h *FindingHandler) CreateFinding(c *gin.Context) {
	var input struct {
		ScanID      string       `json:"scan_id" binding:"required"`
		TargetID    string       `json:"target_id" binding:"required"`
		Title       string       `json:"title" binding:"required"`
		Description string       `json:"description"`
		Severity    string       `json:"severity" binding:"required"`
		FindingType string       `json:"finding_type" binding:"required"`
		Details     models.JSONB `json:"details"`
		Verified    bool         `json:"verified"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate severity
	if input.Severity != models.SeverityCritical &&
		input.Severity != models.SeverityHigh &&
		input.Severity != models.SeverityMedium &&
		input.Severity != models.SeverityLow &&
		input.Severity != models.SeverityInfo {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid severity level"})
		return
	}

	scanID, err := uuid.Parse(input.ScanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan ID format"})
		return
	}

	targetID, err := uuid.Parse(input.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	finding := &models.Finding{
		ScanID:      scanID,
		TargetID:    targetID,
		Title:       input.Title,
		Description: input.Description,
		Severity:    input.Severity,
		FindingType: input.FindingType,
		Details:     input.Details,
		Verified:    input.Verified,
		Fixed:       false, // Default to not fixed
	}

	err = h.findingService.Create(finding)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create finding"})
		return
	}

	c.JSON(http.StatusCreated, finding)
}

// UpdateFinding updates an existing finding
// @Summary Update a finding
// @Description Update an existing finding
// @Tags findings
// @Accept json
// @Produce json
// @Param id path string true "Finding ID"
// @Param finding body object true "Updated Finding Details"
// @Success 200 {object} models.Finding
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/findings/{id} [put]
func (h *FindingHandler) UpdateFinding(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid finding ID format"})
		return
	}

	var input struct {
		Title       string       `json:"title"`
		Description string       `json:"description"`
		Severity    string       `json:"severity"`
		FindingType string       `json:"finding_type"`
		Details     models.JSONB `json:"details"`
		Verified    *bool        `json:"verified"`
		Fixed       *bool        `json:"fixed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	finding, err := h.findingService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Finding not found"})
		return
	}

	// Update fields if provided
	if input.Title != "" {
		finding.Title = input.Title
	}

	if input.Description != "" {
		finding.Description = input.Description
	}

	if input.Severity != "" {
		// Validate severity
		if input.Severity != models.SeverityCritical &&
			input.Severity != models.SeverityHigh &&
			input.Severity != models.SeverityMedium &&
			input.Severity != models.SeverityLow &&
			input.Severity != models.SeverityInfo {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid severity level"})
			return
		}
		finding.Severity = input.Severity
	}

	if input.FindingType != "" {
		finding.FindingType = input.FindingType
	}

	if input.Details != nil {
		finding.Details = input.Details
	}

	if input.Verified != nil {
		finding.Verified = *input.Verified
	}

	if input.Fixed != nil {
		finding.Fixed = *input.Fixed
	}

	err = h.findingService.Update(finding)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update finding"})
		return
	}

	c.JSON(http.StatusOK, finding)
}

// DeleteFinding deletes a finding
// @Summary Delete a finding
// @Description Delete a finding
// @Tags findings
// @Accept json
// @Produce json
// @Param id path string true "Finding ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/findings/{id} [delete]
func (h *FindingHandler) DeleteFinding(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid finding ID format"})
		return
	}

	err = h.findingService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete finding"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Finding deleted successfully"})
}

// BulkUpdateFindings updates multiple findings (e.g., mark multiple as fixed)
// @Summary Bulk update findings
// @Description Update multiple findings at once
// @Tags findings
// @Accept json
// @Produce json
// @Param bulkUpdate body object true "Bulk Update Details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/findings/bulk-update [post]
func (h *FindingHandler) BulkUpdateFindings(c *gin.Context) {
	var input struct {
		FindingIDs []string `json:"finding_ids" binding:"required"`
		Fixed      *bool    `json:"fixed"`
		Verified   *bool    `json:"verified"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Fixed == nil && input.Verified == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No update parameters provided"})
		return
	}

	// Convert string IDs to UUID
	var findingIDs []uuid.UUID
	for _, idStr := range input.FindingIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid finding ID format: " + idStr})
			return
		}
		findingIDs = append(findingIDs, id)
	}

	err := h.findingService.BulkUpdate(findingIDs, input.Fixed, input.Verified)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update findings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Findings updated successfully"})
}

// MarkFixed marks a finding as fixed
// @Summary Mark a finding as fixed
// @Description Mark a finding as fixed or not fixed
// @Tags findings
// @Accept json
// @Produce json
// @Param id path string true "Finding ID"
// @Param fixed path boolean true "Fixed status"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/findings/{id}/fixed/{fixed} [put]
func (h *FindingHandler) MarkFixed(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid finding ID format"})
		return
	}

	fixedStr := c.Param("fixed")
	var fixed bool
	if fixedStr == "true" {
		fixed = true
	} else if fixedStr == "false" {
		fixed = false
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fixed parameter must be 'true' or 'false'"})
		return
	}

	err = h.findingService.MarkFixed(id, fixed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update finding"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Finding updated successfully"})
}

// MarkVerified marks a finding as verified
// @Summary Mark a finding as verified
// @Description Mark a finding as verified or not verified
// @Tags findings
// @Accept json
// @Produce json
// @Param id path string true "Finding ID"
// @Param verified path boolean true "Verified status"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/findings/{id}/verified/{verified} [put]
func (h *FindingHandler) MarkVerified(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid finding ID format"})
		return
	}

	verifiedStr := c.Param("verified")
	var verified bool
	if verifiedStr == "true" {
		verified = true
	} else if verifiedStr == "false" {
		verified = false
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verified parameter must be 'true' or 'false'"})
		return
	}

	err = h.findingService.MarkVerified(id, verified)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update finding"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Finding updated successfully"})
}

// GetFindingsSummary returns a summary of findings by severity for a project
// @Summary Get findings summary
// @Description Get a summary of findings by severity for a project
// @Tags findings
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Success 200 {object} map[string]int64
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/findings/summary/{project_id} [get]
func (h *FindingHandler) GetFindingsSummary(c *gin.Context) {
	projectIDStr := c.Param("project_id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	summary, err := h.findingService.CountBySeverity(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve findings summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetLatestFindings returns the most recent findings
// @Summary Get latest findings
// @Description Get the most recent findings
// @Tags findings
// @Accept json
// @Produce json
// @Param limit query int false "Number of findings to return (default 10)"
// @Success 200 {array} models.Finding
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/findings/latest [get]
func (h *FindingHandler) GetLatestFindings(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit := 10
	// Parse limit if provided
	if limitStr != "10" {
		if _, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}

	findings, err := h.findingService.GetLatestFindings(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve latest findings"})
		return
	}

	c.JSON(http.StatusOK, findings)
}
