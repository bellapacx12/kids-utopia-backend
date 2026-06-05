package handler

import (
	"net/http"

	"github.com/bellapacx/kids-utopia/internal/gamification/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func New(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetXP(c *gin.Context) {

	childID := c.Param("childId")

	xp, err := h.service.GetXP(c.Request.Context(), childID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, xp)
}