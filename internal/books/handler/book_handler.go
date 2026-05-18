package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/books/dto"
	"github.com/bellapacx/kids-utopia/internal/books/service"
)

type BookHandler struct {
	service *service.BookService
}

func NewBookHandler(s *service.BookService) *BookHandler {
	return &BookHandler{service: s}
}
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req dto.CreateBookRequest

	// 1. Bind request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// 2. Call service
	book, err := h.service.CreateBook(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create book",
		})
		return
	}

	// 3. Response
	c.JSON(http.StatusCreated, gin.H{
		"data": book,
	})
}
func (h *BookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")

	book, err := h.service.GetBookByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "book not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})
}
func (h *BookHandler) UploadBook(c *gin.Context) {

	log.Println("📥 [UploadBook] endpoint hit")

	file, err := c.FormFile("file")
	if err != nil {
		log.Println("❌ [UploadBook] file missing:", err)
		c.JSON(400, gin.H{"error": "file required"})
		return
	}

	log.Println("📄 [UploadBook] file received:", file.Filename)

	src, err := file.Open()
	if err != nil {
		log.Println("❌ [UploadBook] cannot open file:", err)
		c.JSON(500, gin.H{"error": "cannot open file"})
		return
	}
	defer src.Close()

	log.Println("⬆️ [UploadBook] uploading to storage...")

	// upload to S3
	fileURL, err := h.service.UploadToStorage(c.Request.Context(), src, file.Filename)
	if err != nil {
		log.Println("❌ [UploadBook] upload failed:", err)
		c.JSON(500, gin.H{"error": "upload failed"})
		return
	}

	log.Println("✅ [UploadBook] upload success:", fileURL)

	log.Println("📘 [UploadBook] creating uploaded book...")

	book, err := h.service.CreateUploadedBook(c.Request.Context(), file.Filename, fileURL)
	if err != nil {
		log.Println("❌ [UploadBook] db error:", err)
		c.JSON(500, gin.H{"error": "db errorrrr"})
		return
	}

	log.Println("🎉 [UploadBook] book created:", book.ID)

	c.JSON(200, gin.H{
		"data": book,
	})
}