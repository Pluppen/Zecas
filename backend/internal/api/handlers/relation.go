package handlers

import (
	"net/http"

	"backend/internal/models"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RelationHandler struct {
	relationService *services.RelationService
	targetService   *services.TargetService
}

func NewRelationHandler(relationService *services.RelationService, targetService *services.TargetService) *RelationHandler {
	return &RelationHandler{
		relationService: relationService,
		targetService:   targetService,
	}
}

// GetRelations returns all target relations
// @Summary Get target relations
// @Description Get all target relations with optional filtering
// @Tags relations
// @Accept json
// @Produce json
// @Param source_id query string false "Source Target ID"
// @Param destination_id query string false "Destination Target ID"
// @Param relation_type query string false "Relation Type"
// @Success 200 {array} models.TargetRelation
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/relations [get]
func (h *RelationHandler) GetRelations(c *gin.Context) {
	// Get query parameters
	sourceIDStr := c.Query("source_id")
	destinationIDStr := c.Query("destination_id")
	relationType := c.Query("relation_type")

	var sourceID, destinationID *uuid.UUID
	var err error

	// Parse UUIDs if provided
	if sourceIDStr != "" {
		id, err := uuid.Parse(sourceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source ID format"})
			return
		}
		sourceID = &id
	}

	if destinationIDStr != "" {
		id, err := uuid.Parse(destinationIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid destination ID format"})
			return
		}
		destinationID = &id
	}

	// Get relations with optional filtering
	relations, err := h.relationService.GetFiltered(sourceID, destinationID, relationType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve relations"})
		return
	}

	c.JSON(http.StatusOK, relations)
}

// GetRelation returns a specific relation by ID
// @Summary Get a target relation
// @Description Get a specific target relation by ID
// @Tags relations
// @Accept json
// @Produce json
// @Param id path string true "Relation ID"
// @Success 200 {object} models.TargetRelation
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/relations/{id} [get]
func (h *RelationHandler) GetRelation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid relation ID format"})
		return
	}

	relation, err := h.relationService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Relation not found"})
		return
	}

	c.JSON(http.StatusOK, relation)
}

// CreateRelation creates a new relationship between two targets
// @Summary Create a target relation
// @Description Create a new relationship between two targets
// @Tags relations
// @Accept json
// @Produce json
// @Param relation body object true "Relation Details"
// @Success 201 {object} models.TargetRelation
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/relations [post]
func (h *RelationHandler) CreateRelation(c *gin.Context) {
	var input struct {
		SourceID      string       `json:"source_id" binding:"required"`
		DestinationID string       `json:"destination_id" binding:"required"`
		RelationType  string       `json:"relation_type" binding:"required"`
		Metadata      models.JSONB `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse UUIDs
	sourceID, err := uuid.Parse(input.SourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source ID format"})
		return
	}

	destinationID, err := uuid.Parse(input.DestinationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid destination ID format"})
		return
	}

	// Verify both targets exist
	_, err = h.targetService.GetByID(sourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source target not found"})
		return
	}

	_, err = h.targetService.GetByID(destinationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Destination target not found"})
		return
	}

	// Validate relation type
	validRelationTypes := map[string]bool{
		models.RelationResolvesTo:   true,
		models.RelationParentOf:     true,
		models.RelationChildOf:      true,
		models.RelationHostsService: true,
	}

	if !validRelationTypes[input.RelationType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid relation type"})
		return
	}

	// Create the relation
	relation := models.TargetRelation{
		SourceID:      sourceID,
		DestinationID: destinationID,
		RelationType:  input.RelationType,
		Metadata:      input.Metadata,
	}

	// Save the relation
	err = h.relationService.Create(&relation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create relation"})
		return
	}

	c.JSON(http.StatusCreated, relation)
}

// DeleteRelation deletes a relation
// @Summary Delete a target relation
// @Description Delete a relation by ID
// @Tags relations
// @Accept json
// @Produce json
// @Param id path string true "Relation ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/relations/{id} [delete]
func (h *RelationHandler) DeleteRelation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid relation ID format"})
		return
	}

	err = h.relationService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete relation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Relation deleted successfully"})
}

// BulkCreateRelations creates multiple relations at once
// @Summary Bulk create relations
// @Description Create multiple relations at once
// @Tags relations
// @Accept json
// @Produce json
// @Param relations body object true "Relations Details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/relations/bulk [post]
func (h *RelationHandler) BulkCreateRelations(c *gin.Context) {
	var input struct {
		Relations []struct {
			SourceID      string       `json:"source_id" binding:"required"`
			DestinationID string       `json:"destination_id" binding:"required"`
			RelationType  string       `json:"relation_type" binding:"required"`
			Metadata      models.JSONB `json:"metadata"`
		} `json:"relations" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var relations []models.TargetRelation
	validRelationTypes := map[string]bool{
		models.RelationResolvesTo:   true,
		models.RelationParentOf:     true,
		models.RelationChildOf:      true,
		models.RelationHostsService: true,
	}

	// Create and validate each relation
	for _, r := range input.Relations {
		// Validate relation type
		if !validRelationTypes[r.RelationType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid relation type: " + r.RelationType})
			return
		}

		// Parse UUIDs
		sourceID, err := uuid.Parse(r.SourceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source ID format: " + r.SourceID})
			return
		}

		destinationID, err := uuid.Parse(r.DestinationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid destination ID format: " + r.DestinationID})
			return
		}

		// Create relation
		relation := models.TargetRelation{
			SourceID:      sourceID,
			DestinationID: destinationID,
			RelationType:  r.RelationType,
			Metadata:      r.Metadata,
		}

		relations = append(relations, relation)
	}

	// Bulk create relations
	err := h.relationService.BulkCreate(relations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create relations"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully created relations",
		"count":   len(relations),
	})
}

// GetTargetRelations gets all relations for a specific target
// @Summary Get relations for a target
// @Description Get all relations where a specific target is either source or destination
// @Tags relations
// @Accept json
// @Produce json
// @Param target_id path string true "Target ID"
// @Success 200 {array} models.TargetRelation
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/targets/{id}/relations [get]
func (h *RelationHandler) GetTargetRelations(c *gin.Context) {
	targetIDStr := c.Param("id")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID format"})
		return
	}

	relations, err := h.relationService.GetRelationsForTarget(targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve relations"})
		return
	}

	c.JSON(http.StatusOK, relations)
}
