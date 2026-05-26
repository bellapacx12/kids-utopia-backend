package library

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

// =========================
// CONTINUE READING
// =========================

func (h *Handler) GetContinueReading(c *gin.Context) {

	childID := c.Param("childId")

	if childID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing child id",
		})
		return
	}

	response, err := h.service.GetContinueReading(
		c.Request.Context(),
		childID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}