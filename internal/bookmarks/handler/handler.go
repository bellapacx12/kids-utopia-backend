package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/bookmarks/model"
	"github.com/bellapacx/kids-utopia/internal/bookmarks/service"
	"github.com/bellapacx/kids-utopia/pkg/contextkeys"
)
type Handler struct {
	service *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{service: s}
}
func (h *Handler) Create(c *gin.Context) {
	ctx := c.Request.Context()



	var req struct {
		ChildID string `json:"child_id"`
		BookID  string `json:"book_id"`
		Page    int    `json:"page"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	err := h.service.Create(ctx, &model.Bookmark{
		ChildID: req.ChildID,
		BookID:  req.BookID,
		Page:    req.Page,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "bookmark_created",
	})
}
func (h *Handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		ChildID string `json:"child_id"`
		BookID  string `json:"book_id"`
		Page    int    `json:"page"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	err := h.service.Delete(ctx, req.ChildID, req.BookID, req.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "bookmark_deleted",
	})
}
func (h *Handler) ListByBook(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		ChildID string `json:"child_id"`
		BookID  string `json:"book_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	result, err := h.service.ListByBook(ctx, req.ChildID, req.BookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
func (h *Handler) ListByChild(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		ChildID string `json:"child_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	result, err := h.service.ListByChild(ctx, req.ChildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
func (h *Handler) ListDetailedByChild(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetString(contextkeys.UserID)
    
	
	var req struct {
		ChildID string `json:"child_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// ✅ REQUIRED VALIDATION (THIS IS WHAT YOU ARE MISSING)
	if req.ChildID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "child_id is required",
		})
		return
	}

	result, err := h.service.ListDetailedByChild(ctx, userID, req.ChildID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}