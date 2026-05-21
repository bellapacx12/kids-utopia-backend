package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/access/service"
	"github.com/bellapacx/kids-utopia/internal/books/repository"
)

type Middleware struct {
	accessService *service.Service
	bookRepo      repository.BookRepository
}

func New(a *service.Service, b repository.BookRepository) *Middleware {
	return &Middleware{
		accessService: a,
		bookRepo:      b,
	}
}
func (m *Middleware) CheckBookAccess() gin.HandlerFunc {
	return func(c *gin.Context) {

		bookID := c.Param("id")
		userID := c.GetString("user_id") // empty if guest

		book, err := m.bookRepo.FindByID(c, bookID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			c.Abort()
			return
		}

		ok, err := m.accessService.CanAccessBook(c, userID, book)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "access error"})
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