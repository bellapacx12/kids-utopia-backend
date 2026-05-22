package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/subscriptions/handler"
	"github.com/bellapacx/kids-utopia/pkg/config"
	"github.com/bellapacx/kids-utopia/pkg/middleware"
)

func RegisterSubscriptionRoutes(
	rg *gin.RouterGroup,
	h *handler.Handler,
) {
	cfg := config.Load()
	protected := rg.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret)) // 👈 ADD THIS

	protected.POST("/start", h.Start)
	protected.GET("/me", h.Me)
}