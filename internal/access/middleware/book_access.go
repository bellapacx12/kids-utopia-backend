package middleware

import (
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
		// NORMAL ACCESS FLOW
		// =========================

		bookID := c.Param("id")

		userID := c.GetString(contextkeys.UserID)

		book, err := m.bookRepo.FindByID(c, bookID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "book not found",
			})

			c.Abort()
			return
		}

		ok, err := m.accessService.CanAccessBook(
			c,
			userID,
			book,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "access error",
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