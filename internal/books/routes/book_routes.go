package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/books/handler"
)

// =========================
// READER ROUTES
// =========================

func RegisterReaderRoutes(
	books *gin.RouterGroup,
	h *handler.BookHandler,
) {

	books.GET("/", h.ListBooks)

	books.GET("/:id", h.GetBooks)
}

// =========================
// EDITOR ROUTES
// =========================

func RegisterEditorBookRoutes(
	books *gin.RouterGroup,
	h *handler.BookHandler,
) {

	books.POST("/", h.CreateBook)

	books.POST("/upload", h.UploadBook)
}