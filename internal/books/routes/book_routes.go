package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/books/handler"
)

func RegisterBookRoutes(r *gin.Engine, h *handler.BookHandler) {
	api := r.Group("/api/v1")
	books := api.Group("/books")

	{
		books.POST("/", h.CreateBook)
		books.GET("/:id", h.GetBook)
		books.POST("/upload", h.UploadBook)
	}
}