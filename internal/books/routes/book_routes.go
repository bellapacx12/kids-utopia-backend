package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/books/handler"
)

func RegisterBookRoutes(
	books *gin.RouterGroup,
	h *handler.BookHandler,
) {

	books.POST("/", h.CreateBook)

	books.GET("/:id", h.GetBook)

	books.POST("/upload", h.UploadBook)
}