package users

import (
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/users/handler"
)
func RegisterRoutes(
	r *gin.RouterGroup,
	h *handler.UserHandler,
	authMiddleware gin.HandlerFunc,
) {

	users := r.Group("/users")
	users.Use(authMiddleware)

	users.GET("/me", h.Me)
	users.PATCH("/me", h.UpdateMe)
}