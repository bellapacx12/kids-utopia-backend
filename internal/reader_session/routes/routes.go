package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/reader_session/handler"
)

func RegisterRoutes(
	rg *gin.RouterGroup,
	h *handler.Handler,
) {

	sessions := rg.Group("/reader/sessions")

	// =========================
	// SESSION LIFECYCLE
	// =========================

	sessions.POST("/start", h.StartSession)
	sessions.POST("/end", h.EndSession)

	// =========================
	// SESSION STATE (IMPORTANT FOR READER UX)
	// =========================

	sessions.GET("/:id", h.GetSession)          // fetch session details
	 // current running session
   // progress update (page changes)
}