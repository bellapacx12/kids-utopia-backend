package children

import (
	"github.com/gin-gonic/gin"

	analyticshandler "github.com/bellapacx/kids-utopia/internal/analytics/handler"
	"github.com/bellapacx/kids-utopia/internal/children/handler"
	streakhandler "github.com/bellapacx/kids-utopia/internal/streak/handler"
)

func RegisterRoutes(
	r *gin.RouterGroup,
	h *handler.ChildHandler,
	auth gin.HandlerFunc,
	streak *streakhandler.Handler,
	analytics *analyticshandler.Handler,
) {

	children := r.Group("/children")
	children.Use(auth)

	children.POST("/", h.Create)
	children.GET("/", h.MyChildren)
	children.GET("/:childId/streak", streak.Get)
children.GET("/:childId/analytics", analytics.GetAnalytics)
}