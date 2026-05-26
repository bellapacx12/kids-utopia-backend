package bootstrap

import "github.com/gin-gonic/gin"

func (a *App) registerHealth() {

	a.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}