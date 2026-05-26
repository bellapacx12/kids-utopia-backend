package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/reader/engine"
)

func RegisterRoutes(
	rg *gin.RouterGroup,
	e *engine.Engine,
) {

	reader := rg.Group("/reader")

	reader.POST("/open", e.OpenHandler)
	reader.POST("/update", e.UpdateHandler)
	reader.POST("/close", e.CloseHandler)
	reader.GET(
	"/state/:bookId/:childId",
	e.StateHandler,
)
}