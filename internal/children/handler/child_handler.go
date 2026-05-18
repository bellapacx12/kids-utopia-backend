package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/children/dto"
	"github.com/bellapacx/kids-utopia/internal/children/service"
)

type ChildHandler struct {
	service *service.ChildService
}

func NewChildHandler(s *service.ChildService) *ChildHandler {
	return &ChildHandler{service: s}
}
func (h *ChildHandler) Create(c *gin.Context) {

	parentID := c.GetString("userID")

	var req dto.CreateChildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.Create(c.Request.Context(), parentID, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "child created"})
}
func (h *ChildHandler) MyChildren(c *gin.Context) {

	parentID := c.GetString("userID")

	children, err := h.service.GetByParent(c.Request.Context(), parentID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": children})
}