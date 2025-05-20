package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CertificateHandler struct {
	certificateService *services.CertificateService
}

func NewCertificateHandler(certificateService *services.CertificateService) *CertificateHandler {
	return &CertificateHandler{
		certificateService: certificateService,
	}
}

// GetCertificates returns all certificates with optional filtering
// @Summary Get all certificates
// @Description Get all certificates with optional filtering
// @Tags certificates
// @Accept json
// @Produce json
// @Param severity query string false "Filter by severity"
// @Param type query string false "Filter by certificate type"
// @Param fixed query string false "Filter by fixed status"
// @Param project_id query string false "Filter by project ID"
// @Success 200 {array} models.Certificate
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/certificates [get]
func (h *CertificateHandler) GetCertificates(c *gin.Context) {
	// Support filtering
	certificateType := c.Query("type")

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

	certificates, err := h.certificateService.GetFiltered(projectID, certificateType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve certificates"})
		return
	}

	c.JSON(http.StatusOK, certificates)
}

// GetCertificate returns a specific certificate by ID
// @Summary Get a certificate
// @Description Get a specific certificate by ID
// @Tags certificates
// @Accept json
// @Produce json
// @Param id path string true "Certificate ID"
// @Success 200 {object} models.Certificate
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/certificates/{id} [get]
func (h *CertificateHandler) GetCertificate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certificate ID format"})
		return
	}

	certificate, err := h.certificateService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}

	c.JSON(http.StatusOK, certificate)
}

// CreateCertificate creates a new certificate
// @Summary Create a certificate
// @Description Create a new certificate
// @Tags certificates
// @Accept json
// @Produce json
// @Param certificate body object true "Certificate Details"
// @Success 201 {object} models.Certificate
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/certificates [post]
func (h *CertificateHandler) CreateCertificate(c *gin.Context) {
	var input struct {
		ScanID        *string      `json:"scan_id"`
		ProjectID     string       `json:"project_id" binding:"required"`
		TargetID      *string      `json:"target_id"`
		ServiceID     *string      `json:"service_id"`
		ApplicationID *string      `json:"application_id"`
		ExpiresAt     time.Time    `json:"expires_at"`
		IssuedAt      time.Time    `json:"issued_at"`
		Issuer        string       `json:"issuer" binding:"required"`
		Domain        string       `json:"domain" binding:"required"`
		Details       models.JSONB `json:"details"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	var err error
	var targetID uuid.UUID
	var serviceID uuid.UUID
	var projectID uuid.UUID
	var applicationID uuid.UUID

	if input.TargetID != nil {
		targetID, err = uuid.Parse(*input.TargetID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
			return
		}
	}

	if input.ServiceID != nil {
		serviceID, err = uuid.Parse(*input.ServiceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID format"})
			return
		}
	}

	if input.ApplicationID != nil {
		applicationID, err = uuid.Parse(*input.ApplicationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID format"})
			return
		}
	}

	projectID, err = uuid.Parse(input.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	certificate := &models.Certificate{
		ScanID:        scanID,
		ProjectID:     projectID,
		TargetID:      targetID,
		ServiceID:     serviceID,
		ApplicationID: applicationID,
		Details:       input.Details,
		Issuer:        input.Issuer,
		Domain:        input.Domain,
		ExpiresAt:     input.ExpiresAt,
		IssuedAt:      input.IssuedAt,
	}

	err = h.certificateService.Create(certificate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create certificate"})
		return
	}

	c.JSON(http.StatusCreated, certificate)
}

// UpdateCertificate updates an existing certificate
// @Summary Update a certificate
// @Description Update an existing certificate
// @Tags certificates
// @Accept json
// @Produce json
// @Param id path string true "Certificate ID"
// @Param certificate body object true "Updated Certificate Details"
// @Success 200 {object} models.Certificate
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/certificates/{id} [put]
func (h *CertificateHandler) UpdateCertificate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certificate ID format"})
		return
	}

	var input struct {
		TargetID      *string       `json:"target_id"`
		ServiceID     *string       `json:"service_id"`
		ApplicationID *string       `json:"application_id"`
		ExpiresAt     *time.Time    `json:"expires_at"`
		IssuedAt      *time.Time    `json:"issued_at"`
		Issuer        *string       `json:"issuer" binding:"required"`
		Domain        *string       `json:"domain" binding:"required"`
		Details       *models.JSONB `json:"details"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	certificate, err := h.certificateService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}

	var targetID uuid.UUID
	var serviceID uuid.UUID
	var applicationID uuid.UUID

	if input.TargetID != nil {
		targetID, err = uuid.Parse(*input.TargetID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
			return
		}
		certificate.TargetID = targetID
	}

	if input.ServiceID != nil {
		serviceID, err = uuid.Parse(*input.ServiceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID format"})
			return
		}
		certificate.ServiceID = serviceID
	}

	if input.ApplicationID != nil {
		applicationID, err = uuid.Parse(*input.ApplicationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID format"})
			return
		}
		certificate.ApplicationID = applicationID
	}

	if input.ExpiresAt != nil {
		certificate.ExpiresAt = *input.ExpiresAt
	}

	if input.IssuedAt != nil {
		certificate.IssuedAt = *input.IssuedAt
	}

	if input.Issuer != nil {
		certificate.Issuer = *input.Issuer
	}

	if input.Domain != nil {
		certificate.Domain = *input.Domain
	}

	if input.Details != nil {
		certificate.Details = *input.Details
	}

	err = h.certificateService.Update(certificate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update certificate"})
		return
	}

	c.JSON(http.StatusOK, certificate)
}

// DeleteCertificate deletes a certificate
// @Summary Delete a certificate
// @Description Delete a certificate
// @Tags certificates
// @Accept json
// @Produce json
// @Param id path string true "Certificate ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/certificates/{id} [delete]
func (h *CertificateHandler) DeleteCertificate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid certificate ID format"})
		return
	}

	err = h.certificateService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete certificate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certificate deleted successfully"})
}
