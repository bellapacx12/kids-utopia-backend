package handler

import (
	"net/http"

	"github.com/bellapacx/kids-utopia/internal/books/dto"
	"github.com/bellapacx/kids-utopia/internal/books/service"
	"github.com/gin-gonic/gin"
)

type EditorHandler struct {
	service *service.EditorService
}
type UpdateAccessTypeRequest struct {
	AccessType string `json:"access_type"`
}
func NewEditorHandler(s *service.EditorService) *EditorHandler {
	return &EditorHandler{s}
}
func (h *EditorHandler) GetEditor(c *gin.Context) {

	variantID := c.Param("id")
	res, err := h.service.GetEditor(c.Request.Context(), variantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
func (h *EditorHandler) SaveEditor(c *gin.Context) {

	variantID := c.Param("id")

	var req dto.SaveEditorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SaveEditor(c.Request.Context(), variantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "saved"})
}
func (h *EditorHandler) UpdateAccessType(c *gin.Context) {

	bookID := c.Param("id")

	var req UpdateAccessTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if req.AccessType != "free" && req.AccessType != "premium" {
		c.JSON(400, gin.H{"error": "invalid access_type"})
		return
	}

	err := h.service.UpdateAccessType(
		c.Request.Context(),
		bookID,
		req.AccessType,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to update access type"})
		return
	}

	c.JSON(200, gin.H{
		"message": "access type updated",
	})
}