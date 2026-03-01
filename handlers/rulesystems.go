package handlers

import (
	"net/http"

	"github.com/FourneauxThibaut/CF-Back/internal/auth"
	"github.com/FourneauxThibaut/CF-Back/internal/ruleeditor"
	"github.com/gin-gonic/gin"
)

// RuleSystemHandler wraps the rule editor service for HTTP handlers.
type RuleSystemHandler struct {
	svc *ruleeditor.Service
}

// NewRuleSystemHandler returns a handler that uses the given service.
func NewRuleSystemHandler(svc *ruleeditor.Service) *RuleSystemHandler {
	return &RuleSystemHandler{svc: svc}
}

// ListSystems GET /api/rule-systems
func (h *RuleSystemHandler) ListSystems(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	list, err := h.svc.ListSystems(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetSystem GET /api/rule-systems/:systemId
func (h *RuleSystemHandler) GetSystem(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	sys, err := h.svc.GetSystem(c.Request.Context(), systemID, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sys)
}

// CreateSystem POST /api/rule-systems
func (h *RuleSystemHandler) CreateSystem(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sys, err := h.svc.CreateSystem(c.Request.Context(), userID, body.Name, body.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sys)
}

// UpdateSystem PUT /api/rule-systems/:systemId
func (h *RuleSystemHandler) UpdateSystem(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	var body struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sys, err := h.svc.UpdateSystem(c.Request.Context(), systemID, userID, body.Name, body.Description)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sys)
}

// DeleteSystem DELETE /api/rule-systems/:systemId
func (h *RuleSystemHandler) DeleteSystem(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	if err := h.svc.DeleteSystem(c.Request.Context(), systemID, userID); err != nil {
		if err.Error() == "rule system not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateRule POST /api/rule-systems/:systemId/rules
func (h *RuleSystemHandler) CreateRule(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Order       int    `json:"order"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Name == "" {
		body.Name = "Nouvelle règle"
	}
	rule, err := h.svc.AddRule(c.Request.Context(), systemID, userID, body.Name, body.Description, body.Icon, body.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

// UpdateRule PUT /api/rule-systems/:systemId/rules/:ruleId
func (h *RuleSystemHandler) UpdateRule(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	ruleID := c.Param("ruleId")
	var body struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Icon        *string `json:"icon"`
		IsActive    *bool   `json:"isActive"`
		Order       *int    `json:"order"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rule, err := h.svc.UpdateRule(c.Request.Context(), systemID, ruleID, userID, body.Name, body.Description, body.Icon, body.IsActive, body.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

// DeleteRule DELETE /api/rule-systems/:systemId/rules/:ruleId
func (h *RuleSystemHandler) DeleteRule(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	ruleID := c.Param("ruleId")
	if err := h.svc.DeleteRule(c.Request.Context(), systemID, ruleID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ReorderRules PUT /api/rule-systems/:systemId/rules/reorder
func (h *RuleSystemHandler) ReorderRules(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	var body struct {
		OrderedIDs []string `json:"orderedIds"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ReorderRules(c.Request.Context(), systemID, userID, body.OrderedIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// AddBlock POST /api/rule-systems/:systemId/rules/:ruleId/blocks
func (h *RuleSystemHandler) AddBlock(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	ruleID := c.Param("ruleId")
	var block ruleeditor.RuleBlock
	if err := c.ShouldBindJSON(&block); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := h.svc.AddBlock(c.Request.Context(), systemID, ruleID, userID, block)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// UpdateBlock PUT /api/rule-systems/:systemId/rules/:ruleId/blocks/:blockId
func (h *RuleSystemHandler) UpdateBlock(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	ruleID := c.Param("ruleId")
	blockID := c.Param("blockId")
	var body struct {
		Segments []ruleeditor.Segment `json:"segments"`
		Order    *int                `json:"order"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.svc.UpdateBlock(c.Request.Context(), systemID, ruleID, blockID, userID, body.Segments, body.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteBlock DELETE /api/rule-systems/:systemId/rules/:ruleId/blocks/:blockId
func (h *RuleSystemHandler) DeleteBlock(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	ruleID := c.Param("ruleId")
	blockID := c.Param("blockId")
	if err := h.svc.DeleteBlock(c.Request.Context(), systemID, ruleID, blockID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ReorderBlocks PUT /api/rule-systems/:systemId/rules/:ruleId/blocks/reorder
func (h *RuleSystemHandler) ReorderBlocks(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	ruleID := c.Param("ruleId")
	var body struct {
		OrderedIDs []string `json:"orderedIds"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ReorderBlocks(c.Request.Context(), systemID, ruleID, userID, body.OrderedIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetBlockDefinitions GET /api/rule-systems/:systemId/block-definitions
func (h *RuleSystemHandler) GetBlockDefinitions(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	defs, err := h.svc.GetBlockDefinitions(c.Request.Context(), systemID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, defs)
}

// CreateBlockDefinition POST /api/rule-systems/:systemId/block-definitions
func (h *RuleSystemHandler) CreateBlockDefinition(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	var def ruleeditor.BlockDefinition
	if err := c.ShouldBindJSON(&def); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := h.svc.AddBlockDefinition(c.Request.Context(), systemID, userID, def)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// UpdateBlockDefinition PUT /api/rule-systems/:systemId/block-definitions/:defId
func (h *RuleSystemHandler) UpdateBlockDefinition(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	defID := c.Param("defId")
	var def ruleeditor.BlockDefinition
	if err := c.ShouldBindJSON(&def); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.svc.UpdateBlockDefinition(c.Request.Context(), systemID, defID, userID, def)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DeleteBlockDefinition DELETE /api/rule-systems/:systemId/block-definitions/:defId
func (h *RuleSystemHandler) DeleteBlockDefinition(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	systemID := c.Param("systemId")
	defID := c.Param("defId")
	if err := h.svc.DeleteBlockDefinition(c.Request.Context(), systemID, defID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
