package handler

import (
	"net/http"

	"github.com/bellapacx/kids-utopia/internal/streak/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.StreakService
}

func New(s *service.StreakService) *Handler {
	return &Handler{service: s}
}
func (h *Handler) Get(c *gin.Context) {

	childID := c.Param("childId")

	streak, err := h.service.GetStreak(c.Request.Context(), childID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "streak not found",
		})
		return
	}

	c.JSON(http.StatusOK, streak)
}