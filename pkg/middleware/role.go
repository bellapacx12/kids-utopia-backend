package middleware

import (
	"net/http"

	"github.com/bellapacx/kids-utopia/pkg/auth"
	"github.com/gin-gonic/gin"
)
func RequireRoles(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		role := auth.GetRole(c)

		for _, r := range allowed {
			if role == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
		})
		c.Abort()
	}
}