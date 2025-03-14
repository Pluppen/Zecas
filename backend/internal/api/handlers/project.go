package handlers

import (
	"log"
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectService *services.ProjectService
	targetService  *services.TargetService
}

func NewProjectHandler(projectService *services.ProjectService, targetService *services.TargetService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		targetService:  targetService,
	}
}

// GetProjects returns all projects
// @Summary Get all projects
// @Description Get all projects with basic information
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {array} models.Project
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects [get]
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	projects, err := h.projectService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetProject returns a specific project by ID
// @Summary Get a specific project
// @Description Get detailed information about a project by ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	project, err := h.projectService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// CreateProject creates a new project
// @Summary Create a new project
// @Description Create a new project with targets
// @Tags projects
// @Accept json
// @Produce json
// @Param project body models.CreateProjectInput true "Project Details"
// @Success 201 {object} models.Project
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var input models.CreateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create project
	project := &models.Project{
		Name:        input.Name,
		Description: input.Description,
	}

	err := h.projectService.Create(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	// Add targets
	for _, ip := range input.IPRanges {
		target := &models.Target{
			ProjectID:  project.ID,
			TargetType: models.TargetTypeIP,
			Value:      ip,
		}
		if err := h.targetService.Create(target); err != nil {
			log.Printf("Error creating IP target: %v", err)
		}
	}

	for _, cidr := range input.CIDRRanges {
		target := &models.Target{
			ProjectID:  project.ID,
			TargetType: models.TargetTypeCIDR,
			Value:      cidr,
		}
		if err := h.targetService.Create(target); err != nil {
			log.Printf("Error creating CIDR target: %v", err)
		}
	}

	for _, domain := range input.Domains {
		target := &models.Target{
			ProjectID:  project.ID,
			TargetType: models.TargetTypeDomain,
			Value:      domain,
		}
		if err := h.targetService.Create(target); err != nil {
			log.Printf("Error creating domain target: %v", err)
		}
	}

	// Get fresh project with targets
	projectWithTargets, err := h.projectService.GetByID(project.ID)
	if err != nil {
		c.JSON(http.StatusCreated, project)
		return
	}

	c.JSON(http.StatusCreated, projectWithTargets)
}

// UpdateProject updates an existing project
// @Summary Update a project
// @Description Update an existing project's details
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param project body models.Project true "Updated Project Details"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.projectService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	project.Name = input.Name
	project.Description = input.Description

	err = h.projectService.Update(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject deletes a project
// @Summary Delete a project
// @Description Delete a project and all associated data
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	err = h.projectService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// GetProjectTargets gets all targets for a project
// @Summary Get project targets
// @Description Get all targets for a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {array} models.Target
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects/{id}/targets [get]
func (h *ProjectHandler) GetProjectTargets(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	targets, err := h.targetService.GetByProjectID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve targets"})
		return
	}

	c.JSON(http.StatusOK, targets)
}

// GetProjectScans gets all scans for a project
// @Summary Get project scans
// @Description Get all scans for a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {array} models.Scan
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects/{id}/scans [get]
func (h *ProjectHandler) GetProjectScans(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	scans, err := h.projectService.GetScans(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scans"})
		return
	}

	c.JSON(http.StatusOK, scans)
}

// GetProjectFindings gets all findings for a project
// @Summary Get project findings
// @Description Get all findings for a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {array} models.Finding
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects/{id}/findings [get]
func (h *ProjectHandler) GetProjectFindings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	findings, err := h.projectService.GetFindings(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve findings"})
		return
	}

	c.JSON(http.StatusOK, findings)
}

// GetProjectApplications gets all findings for a project
// @Summary Get project applications
// @Description Get all applications for a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {array} models.Finding
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects/{id}/applications [get]
func (h *ProjectHandler) GetProjectApplications(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	applications, err := h.projectService.GetApplications(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
		return
	}

	c.JSON(http.StatusOK, applications)
}

// GetProjectServices gets all services for a project
// @Summary Get project services
// @Description Get all services for a specific project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {array} models.Services
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/projects/{id}/services [get]
func (h *ProjectHandler) GetProjectServices(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	services, err := h.projectService.GetServices(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve services"})
		return
	}

	c.JSON(http.StatusOK, services)
}
