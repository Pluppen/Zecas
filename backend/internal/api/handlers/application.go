// internal/api/handlers/application.go
package handlers

import (
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApplicationHandler struct {
	applicationService *services.ApplicationService
	projectService     *services.ProjectService
	targetService      *services.TargetService
	serviceService     *services.ServiceService
}

func NewApplicationHandler(
	applicationService *services.ApplicationService,
	projectService *services.ProjectService,
	targetService *services.TargetService,
	serviceService *services.ServiceService,
) *ApplicationHandler {
	return &ApplicationHandler{
		applicationService: applicationService,
		projectService:     projectService,
		targetService:      targetService,
		serviceService:     serviceService,
	}
}

// GetApplications returns all applications with optional filtering
// @Summary Get all applications
// @Description Get all applications with optional filtering
// @Tags applications
// @Accept json
// @Produce json
// @Param project_id query string false "Filter by project ID"
// @Param type query string false "Filter by application type"
// @Param target_id query string false "Filter by host target ID"
// @Param service_id query string false "Filter by service ID"
// @Param search query string false "Search by name or URL"
// @Success 200 {array} models.Application
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/applications [get]
func (h *ApplicationHandler) GetApplications(c *gin.Context) {
	// Check for query parameters
	projectIDStr := c.Query("project_id")
	appType := c.Query("type")
	targetIDStr := c.Query("target_id")
	serviceIDStr := c.Query("service_id")
	search := c.Query("search")

	var applications []models.Application
	var err error

	// Apply filters based on provided parameters
	if search != "" {
		applications, err = h.applicationService.Search(search)
	} else if projectIDStr != "" {
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
			return
		}
		applications, err = h.applicationService.GetByProjectID(projectID)
	} else if appType != "" {
		applications, err = h.applicationService.GetByType(appType)
	} else if targetIDStr != "" {
		targetID, err := uuid.Parse(targetIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
			return
		}
		applications, err = h.applicationService.GetByTargetID(targetID)
	} else if serviceIDStr != "" {
		serviceID, err := uuid.Parse(serviceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID format"})
			return
		}
		applications, err = h.applicationService.GetByServiceID(serviceID)
	} else {
		applications, err = h.applicationService.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
		return
	}

	c.JSON(http.StatusOK, applications)
}

// GetApplication returns a specific application by ID
// @Summary Get an application
// @Description Get a specific application by ID
// @Tags applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} models.Application
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/applications/{id} [get]
func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID format"})
		return
	}

	application, err := h.applicationService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	c.JSON(http.StatusOK, application)
}

// CreateApplication creates a new application
// @Summary Create an application
// @Description Create a new application
// @Tags applications
// @Accept json
// @Produce json
// @Param application body object true "Application Details"
// @Success 201 {object} models.Application
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/applications [post]
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var input struct {
		ProjectID   string       `json:"project_id" binding:"required"`
		Name        string       `json:"name" binding:"required"`
		Type        string       `json:"type" binding:"required"`
		Version     string       `json:"version"`
		Description string       `json:"description"`
		URL         string       `json:"url"`
		HostTarget  *string      `json:"host_target,omitempty"`
		ServiceID   *string      `json:"service_id,omitempty"`
		Metadata    models.JSONB `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse project ID
	projectID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	// Verify project exists
	_, err = h.projectService.GetByID(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Parse host target ID if provided
	var hostTarget *uuid.UUID
	if input.HostTarget != nil {
		parsedID, err := uuid.Parse(*input.HostTarget)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host target ID format"})
			return
		}

		// Verify target exists
		_, err = h.targetService.GetByID(parsedID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
			return
		}
		hostTarget = &parsedID
	}

	// Parse service ID if provided
	var serviceID *uuid.UUID
	if input.ServiceID != nil {
		parsedID, err := uuid.Parse(*input.ServiceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID format"})
			return
		}

		// Verify service exists
		_, err = h.serviceService.GetByID(parsedID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		serviceID = &parsedID
	}

	// Create the application
	application := &models.Application{
		ProjectID:   projectID,
		Name:        input.Name,
		Type:        input.Type,
		Version:     input.Version,
		Description: input.Description,
		URL:         input.URL,
		HostTarget:  hostTarget,
		ServiceID:   serviceID,
		Metadata:    input.Metadata,
	}

	err = h.applicationService.Create(application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, application)
}

// BulkCreateApplications creates multiple applications
// @Summary Bulk create applications
// @Description Create multiple applications at once
// @Tags applications
// @Accept json
// @Produce json
// @Param applications body object true "Applications Details"
// @Success 201 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/applications/bulk [post]
func (h *ApplicationHandler) BulkCreateApplications(c *gin.Context) {
	var input struct {
		ProjectID    string               `json:"project_id" binding:"required"`
		Applications []models.Application `json:"applications" binding:"required"`
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

	// Verify project exists
	_, err = h.projectService.GetByID(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Set project ID for all applications
	for i := range input.Applications {
		input.Applications[i].ProjectID = projectID
	}

	err = h.applicationService.BulkCreate(input.Applications)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create applications"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully created applications",
		"count":   len(input.Applications),
	})
}

// UpdateApplication updates an existing application
// @Summary Update an application
// @Description Update an existing application
// @Tags applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param application body object true "Updated Application Details"
// @Success 200 {object} models.Application
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/applications/{id} [put]
func (h *ApplicationHandler) UpdateApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID format"})
		return
	}

	var input struct {
		Name        string       `json:"name"`
		Type        string       `json:"type"`
		Version     string       `json:"version"`
		Description string       `json:"description"`
		URL         string       `json:"url"`
		HostTarget  *string      `json:"host_target,omitempty"`
		ServiceID   *string      `json:"service_id,omitempty"`
		Metadata    models.JSONB `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application, err := h.applicationService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Update fields if provided
	if input.Name != "" {
		application.Name = input.Name
	}
	if input.Type != "" {
		application.Type = input.Type
	}
	if input.Version != "" {
		application.Version = input.Version
	}
	if input.Description != "" {
		application.Description = input.Description
	}
	if input.URL != "" {
		application.URL = input.URL
	}
	if input.Metadata != nil {
		application.Metadata = input.Metadata
	}

	// Parse host target ID if provided
	if input.HostTarget != nil {
		if *input.HostTarget == "" {
			application.HostTarget = nil
		} else {
			parsedID, err := uuid.Parse(*input.HostTarget)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host target ID format"})
				return
			}

			// Verify target exists
			_, err = h.targetService.GetByID(parsedID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
				return
			}
			application.HostTarget = &parsedID
		}
	}

	// Parse service ID if provided
	if input.ServiceID != nil {
		if *input.ServiceID == "" {
			application.ServiceID = nil
		} else {
			parsedID, err := uuid.Parse(*input.ServiceID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID format"})
				return
			}

			// Verify service exists
			_, err = h.serviceService.GetByID(parsedID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
				return
			}
			application.ServiceID = &parsedID
		}
	}

	err = h.applicationService.Update(application)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
		return
	}

	c.JSON(http.StatusOK, application)
}

// DeleteApplication deletes an application
// @Summary Delete an application
// @Description Delete an application
// @Tags applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/applications/{id} [delete]
func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID format"})
		return
	}

	err = h.applicationService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application deleted successfully"})
}

// GetApplicationFindings returns all findings for a specific application
// @Summary Get application findings
// @Description Get all findings for a specific application
// @Tags applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {array} models.Finding
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/applications/{id}/findings [get]
func (h *ApplicationHandler) GetApplicationFindings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID format"})
		return
	}

	findings, err := h.applicationService.GetFindings(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve findings"})
		return
	}

	c.JSON(http.StatusOK, findings)
}
