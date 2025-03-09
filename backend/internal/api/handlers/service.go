// internal/api/handlers/service.go
package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ServiceHandler struct {
	serviceService *services.ServiceService
	targetService  *services.TargetService
}

func NewServiceHandler(serviceService *services.ServiceService, targetService *services.TargetService) *ServiceHandler {
	return &ServiceHandler{
		serviceService: serviceService,
		targetService:  targetService,
	}
}

// GetServices returns all services
// @Summary Get all services
// @Description Get all services across all targets
// @Tags services
// @Accept json
// @Produce json
// @Success 200 {array} models.Service
// @Failure 500 {object} map[string]string
// @Router /api/v1/services [get]
func (h *ServiceHandler) GetServices(c *gin.Context) {
	// Check for query parameters
	targetIDStr := c.Query("target_id")
	projectIDStr := c.Query("project_id")
	portStr := c.Query("port")
	serviceName := c.Query("service_name")

	var services []models.Service
	var err error

	// Filter by the provided parameter
	if targetIDStr != "" {
		targetID, err := uuid.Parse(targetIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
			return
		}
		services, err = h.serviceService.GetByTargetID(targetID)
	} else if projectIDStr != "" {
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
			return
		}
		services, err = h.serviceService.GetByProjectID(projectID)
	} else if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid port number"})
			return
		}
		services, err = h.serviceService.GetByPort(port)
	} else if serviceName != "" {
		services, err = h.serviceService.GetByServiceName(serviceName)
	} else {
		services, err = h.serviceService.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
		return
	}

	c.JSON(http.StatusOK, services)
}

// GetService returns a specific service by ID
// @Summary Get a service
// @Description Get a specific service by ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} models.Service
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/services/{id} [get]
func (h *ServiceHandler) GetService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID format"})
		return
	}

	service, err := h.serviceService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, service)
}

// CreateService creates a new service
// @Summary Create a service
// @Description Create a new service for a target
// @Tags services
// @Accept json
// @Produce json
// @Param service body object true "Service Details"
// @Success 201 {object} models.Service
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/services [post]
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var input struct {
		TargetID    string       `json:"target_id" binding:"required"`
		Port        int          `json:"port" binding:"required"`
		Protocol    string       `json:"protocol" binding:"required"`
		ServiceName string       `json:"service_name"`
		Version     string       `json:"version"`
		Title       string       `json:"title"`
		Description string       `json:"description"`
		Banner      string       `json:"banner"`
		RawInfo     models.JSONB `json:"raw_info"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetID, err := uuid.Parse(input.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	// Verify target exists
	target, err := h.targetService.GetByID(targetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	// Create the service
	service := &models.Service{
		TargetID:    target.ID,
		Port:        input.Port,
		Protocol:    input.Protocol,
		ServiceName: input.ServiceName,
		Version:     input.Version,
		Title:       input.Title,
		Description: input.Description,
		Banner:      input.Banner,
		RawInfo:     input.RawInfo,
	}

	err = h.serviceService.Create(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// BulkCreateServices creates multiple services
// @Summary Bulk create services
// @Description Create multiple services at once
// @Tags services
// @Accept json
// @Produce json
// @Param services body object true "Services Details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/services/bulk [post]
func (h *ServiceHandler) BulkCreateServices(c *gin.Context) {
	var input struct {
		TargetID string           `json:"target_id" binding:"required"`
		Services []models.Service `json:"services" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetID, err := uuid.Parse(input.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	// Verify target exists
	_, err = h.targetService.GetByID(targetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		return
	}

	// Set target ID for all services
	for i := range input.Services {
		input.Services[i].TargetID = targetID
	}

	err = h.serviceService.BulkCreate(input.Services)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create services"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully created services",
		"count":   len(input.Services),
	})
}

// UpdateService updates an existing service
// @Summary Update a service
// @Description Update an existing service
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param service body object true "Updated Service Details"
// @Success 200 {object} models.Service
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/services/{id} [put]
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID format"})
		return
	}

	var input struct {
		Port        *int         `json:"port"`
		Protocol    string       `json:"protocol"`
		ServiceName string       `json:"service_name"`
		Version     string       `json:"version"`
		Title       string       `json:"title"`
		Description string       `json:"description"`
		Banner      string       `json:"banner"`
		RawInfo     models.JSONB `json:"raw_info"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.serviceService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	// Update fields if provided
	if input.Port != nil {
		service.Port = *input.Port
	}
	if input.Protocol != "" {
		service.Protocol = input.Protocol
	}
	if input.ServiceName != "" {
		service.ServiceName = input.ServiceName
	}
	if input.Version != "" {
		service.Version = input.Version
	}
	if input.Title != "" {
		service.Title = input.Title
	}
	if input.Description != "" {
		service.Description = input.Description
	}
	if input.Banner != "" {
		service.Banner = input.Banner
	}
	if input.RawInfo != nil {
		service.RawInfo = input.RawInfo
	}

	err = h.serviceService.Update(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, service)
}

// DeleteService deletes a service
// @Summary Delete a service
// @Description Delete a service by ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/services/{id} [delete]
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID format"})
		return
	}

	err = h.serviceService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

// GetServiceFindings returns all findings for a specific service
// @Summary Get service findings
// @Description Get all findings for a specific service
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {array} models.Finding
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/services/{id}/findings [get]
func (h *ServiceHandler) GetServiceFindings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID format"})
		return
	}

	findings, err := h.serviceService.GetFindings(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve findings"})
		return
	}

	c.JSON(http.StatusOK, findings)
}
