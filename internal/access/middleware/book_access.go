package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/access/service"
	"github.com/bellapacx/kids-utopia/internal/books/repository"
	"github.com/bellapacx/kids-utopia/pkg/contextkeys"
)

type Middleware struct {
	accessService *service.Service
	bookRepo      repository.BookRepository
}

func New(
	a *service.Service,
	b repository.BookRepository,
) *Middleware {
	return &Middleware{
		accessService: a,
		bookRepo:      b,
	}
}

func (m *Middleware) CheckBookAccess() gin.HandlerFunc {
	return func(c *gin.Context) {

		// =========================
		// ROLE BYPASS
		// =========================
		role := c.GetString(contextkeys.Role)

		if role == "editor" ||
			role == "admin" ||
			role == "super_admin" {
			c.Next()
			return
		}

		// =========================
		// BOOK ID
		// =========================
		bookID := c.Param("id")
		if bookID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing book id",
			})
			c.Abort()
			return
		}

		// =========================
		// FETCH BOOK
		// =========================
		book, err := m.bookRepo.FindByID(c, bookID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "book not found",
			})
			c.Abort()
			return
		}

		// =========================
		// ⭐ FREE ACCESS CHECK (NEW)
		// =========================
		if book.AccessType == "free" {
			c.Next()
			return
		}

		// =========================
		// AUTH CHECK
		// =========================
		userID := c.GetString(contextkeys.UserID)
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		// =========================
		// SUBSCRIPTION CHECK
		// =========================
		ok, err := m.accessService.CanAccessBook(
			c,
			userID,
			book,
		)

		if err != nil {
	log.Printf("❌ CanAccessBook failed: userID=%s bookID=%s err=%v",
		userID, bookID, err)

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
	c.Abort()
	return
}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "premium subscription required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}