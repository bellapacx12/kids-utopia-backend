package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/subscriptions/handler"
)

func RegisterSubscriptionRoutes(
	rg *gin.RouterGroup,
	h *handler.Handler,
) {

	rg.POST("/start", h.Start)

	rg.GET("/me", h.Me)
}