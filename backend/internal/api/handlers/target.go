package handlers

import (
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TargetHandler struct {
	targetService *services.TargetService
}

func NewTargetHandler(targetService *services.TargetService) *TargetHandler {
	return &TargetHandler{
		targetService: targetService,
	}
}

// GetTargets returns all targets
// @Summary Get all targets
// @Description Get all targets across all projects
// @Tags targets
// @Accept json
// @Produce json
// @Success 200 {array} models.Target
// @Failure 500 {object} map[string]string
// @Router /api/v1/targets [get]
func (h *TargetHandler) GetTargets(c *gin.Context) {
	targets, err := h.targetService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve targets"})
		return
	}

	c.JSON(http.StatusOK, targets)
}

// GetTarget returns a specific target by ID
// @Summary Get a target
// @Description Get a specific target by ID
// @Tags targets
// @Accept json
// @Produce json
// @Param id path string true "Target ID"
// @Success 200 {object} models.Target
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/targets/{id} [get]
func (h *TargetHandler) GetTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	target, err := h.targetService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	c.JSON(http.StatusOK, target)
}

// CreateTarget creates a new target
// @Summary Create a target
// @Description Create a new target for a project
// @Tags targets
// @Accept json
// @Produce json
// @Param target body object true "Target Details"
// @Success 201 {object} models.Target
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/targets [post]
func (h *TargetHandler) CreateTarget(c *gin.Context) {
	var input struct {
		ProjectID  string       `json:"project_id" binding:"required"`
		TargetType string       `json:"target_type" binding:"required"`
		Value      string       `json:"value" binding:"required"`
		Metadata   models.JSONB `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate target type
	if input.TargetType != models.TargetTypeIP &&
		input.TargetType != models.TargetTypeCIDR &&
		input.TargetType != models.TargetTypeDomain {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target type"})
		return
	}

	projectID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	target := &models.Target{
		ProjectID:  projectID,
		TargetType: input.TargetType,
		Value:      input.Value,
		Metadata:   input.Metadata,
	}

	err = h.targetService.Create(target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create target"})
		return
	}

	c.JSON(http.StatusCreated, target)
}

// BulkCreateTargets creates multiple targets
// @Summary Bulk create targets
// @Description Create multiple targets at once
// @Tags targets
// @Accept json
// @Produce json
// @Param targets body object true "Targets Details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/targets/bulk [post]
func (h *TargetHandler) BulkCreateTargets(c *gin.Context) {
	var input struct {
		ProjectID  string       `json:"project_id" binding:"required"`
		IPRanges   []string     `json:"ip_ranges"`
		CIDRRanges []string     `json:"cidr_ranges"`
		Domains    []string     `json:"domains"`
		Metadata   models.JSONB `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// Create all targets
	var targets []models.Target

	for _, ip := range input.IPRanges {
		targets = append(targets, models.Target{
			ProjectID:  projectID,
			TargetType: models.TargetTypeIP,
			Value:      ip,
			Metadata:   input.Metadata,
		})
	}

	for _, cidr := range input.CIDRRanges {
		targets = append(targets, models.Target{
			ProjectID:  projectID,
			TargetType: models.TargetTypeCIDR,
			Value:      cidr,
			Metadata:   input.Metadata,
		})
	}

	for _, domain := range input.Domains {
		targets = append(targets, models.Target{
			ProjectID:  projectID,
			TargetType: models.TargetTypeDomain,
			Value:      domain,
			Metadata:   input.Metadata,
		})
	}

	if len(targets) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No targets provided"})
		return
	}

	if err := h.targetService.BulkCreate(targets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create targets"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully created targets",
		"count":   len(targets),
	})
}

// UpdateTarget updates an existing target
// @Summary Update a target
// @Description Update an existing target
// @Tags targets
// @Accept json
// @Produce json
// @Param id path string true "Target ID"
// @Param target body object true "Updated Target Details"
// @Success 200 {object} models.Target
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/targets/{id} [put]
func (h *TargetHandler) UpdateTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	var input struct {
		TargetType string       `json:"target_type"`
		Value      string       `json:"value"`
		Metadata   models.JSONB `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target, err := h.targetService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	// Only update provided fields
	if input.TargetType != "" {
		if input.TargetType != models.TargetTypeIP &&
			input.TargetType != models.TargetTypeCIDR &&
			input.TargetType != models.TargetTypeDomain {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target type"})
			return
		}
		target.TargetType = input.TargetType
	}

	if input.Value != "" {
		target.Value = input.Value
	}

	if input.Metadata != nil {
		target.Metadata = input.Metadata
	}

	err = h.targetService.Update(target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update target"})
		return
	}

	c.JSON(http.StatusOK, target)
}

// DeleteTarget deletes a target
// @Summary Delete a target
// @Description Delete a target by ID
// @Tags targets
// @Accept json
// @Produce json
// @Param id path string true "Target ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/targets/{id} [delete]
func (h *TargetHandler) DeleteTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	err = h.targetService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete target"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target deleted successfully"})
}

// GetTargetFindings returns all findings for a specific target
// @Summary Get target findings
// @Description Get all findings for a specific target
// @Tags targets
// @Accept json
// @Produce json
// @Param id path string true "Target ID"
// @Success 200 {array} models.Finding
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/targets/{id}/findings [get]
func (h *TargetHandler) GetTargetFindings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	findings, err := h.targetService.GetFindings(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve findings"})
		return
	}

	c.JSON(http.StatusOK, findings)
}
