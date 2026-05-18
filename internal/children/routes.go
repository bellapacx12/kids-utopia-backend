package children

import (
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/children/handler"
)

func RegisterRoutes(
	r *gin.RouterGroup,
	h *handler.ChildHandler,
	auth gin.HandlerFunc,
) {

	children := r.Group("/children")
	children.Use(auth)

	children.POST("/", h.Create)
	children.GET("/", h.MyChildren)
}