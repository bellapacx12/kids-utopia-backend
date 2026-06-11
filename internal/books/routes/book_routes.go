package routes

import (
	"github.com/gin-gonic/gin"

	accessmiddleware "github.com/bellapacx/kids-utopia/internal/access/middleware"
	"github.com/bellapacx/kids-utopia/internal/books/handler"
)

// =========================
// READER ROUTES
// =========================

func RegisterReaderRoutes(
	books *gin.RouterGroup,
	h *handler.BookHandler,
	accessMw *accessmiddleware.Middleware,
) {

	// LIST BOOKS
	books.GET("/", h.ListBook)

	// SINGLE BOOK ACCESS CHECK
	books.GET(
		"/:id",
		accessMw.CheckBookAccess(),
		h.GetBooks,
	)
}

// =========================
// EDITOR ROUTES
// =========================

func RegisterEditorBookRoutes(
	books *gin.RouterGroup,
	h *handler.BookHandler,
) {

	books.POST("/", h.CreateBook)

	books.POST("/upload", h.UploadFirstVariant)
	books.POST("/:id/variants", h.UploadVariant)

}