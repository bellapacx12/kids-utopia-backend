package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/subscriptions/dto"
	"github.com/bellapacx/kids-utopia/internal/subscriptions/service"
)

type Handler struct {
	service *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{service: s}
}
func (h *Handler) Start(c *gin.Context) {

	userID := c.GetString("user_id")

	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.Create(c, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "active",
	})
}
func (h *Handler) Me(c *gin.Context) {

	userID := c.GetString("user_id")

	active, err := h.service.HasActive(c, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed"})
		return
	}

	c.JSON(200, gin.H{
		"active": active,
	})
}