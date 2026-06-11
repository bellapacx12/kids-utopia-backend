package handler

import (
	"net/http"

	"github.com/bellapacx/kids-utopia/internal/recommendation/service"
	"github.com/bellapacx/kids-utopia/pkg/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Recommend(c *gin.Context) {

	ctx := c.Request.Context()

	userID := c.GetString(contextkeys.UserID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	childID := c.Query("child_id")
	if childID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "child_id required",
		})
		return
	}

	result, err := h.service.Recommend(ctx, childID)
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