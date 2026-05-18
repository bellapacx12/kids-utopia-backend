package handler

import (
	"net/http"

	"github.com/bellapacx/kids-utopia/internal/progress/dto"
	"github.com/bellapacx/kids-utopia/internal/progress/service"
	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	service *service.ProgressService
}

func NewProgressHandler(s *service.ProgressService) *ProgressHandler {
	return &ProgressHandler{service: s}
}

func (h *ProgressHandler) Update(c *gin.Context) {

	var req dto.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateProgress(
		c.Request.Context(),
		req.ChildID,
		req.BookID,
		req.Page,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *ProgressHandler) Get(c *gin.Context) {

	childID := c.Param("childId")
	bookID := c.Param("bookId")

	data, err := h.service.GetProgress(c.Request.Context(), childID, bookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}