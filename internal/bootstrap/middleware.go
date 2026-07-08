package bootstrap

import (
	"log"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (a *App) registerMiddlewares() {
	a.Router.RedirectTrailingSlash = false

	// Log the origin for debugging
	a.Router.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		log.Printf("Request Origin: %s, Path: %s", origin, c.Request.URL.Path)
		c.Next()
	})

	a.Router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			log.Printf("Checking origin: %s", origin)
			allowed := strings.HasPrefix(origin, "http://localhost") || origin == ""
			log.Printf("Origin allowed: %v", allowed)
			return allowed
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	a.Router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})
}