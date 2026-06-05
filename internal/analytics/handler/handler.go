package handler

import (
	"net/http"

	"github.com/bellapacx/kids-utopia/internal/analytics/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetAnalytics(c *gin.Context) {

	childID := c.Param("childId")

	data, err := h.svc.GetAnalyticss(c.Request.Context(), childID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}