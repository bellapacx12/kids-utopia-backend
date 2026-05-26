package routes

import (
	"github.com/bellapacx/kids-utopia/internal/progress/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.RouterGroup,
	h *handler.ProgressHandler,
	auth gin.HandlerFunc,
) {

	progress := r.Group("/progress")
	progress.Use(auth)
	progress.GET("/:childId/:bookId", h.Get)
}