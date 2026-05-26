package library

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	rg *gin.RouterGroup,
	h *Handler,
) {

	reader := rg.Group("/reader/library")

	reader.GET(
		"/continue/:childId",
		h.GetContinueReading,
	)
}