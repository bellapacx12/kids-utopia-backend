package routes

import (
	"github.com/bellapacx/kids-utopia/internal/books/handler"
	"github.com/gin-gonic/gin"
)

func RegisterEditorRoutes(rg *gin.RouterGroup, h *handler.EditorHandler) {

	rg.GET("/books/:id/editor", h.GetEditor)
	rg.POST("/books/:id/editor/save", h.SaveEditor)
		rg.PUT("/books/:id/editor/access-type", h.UpdateAccessType)
}