package middleware

import "github.com/gin-gonic/gin"

func RoleGuard(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		role := c.GetString("role")

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatus(403)
	}
}