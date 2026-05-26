package bootstrap

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (a *App) registerMiddlewares() {

	a.Router.RedirectTrailingSlash = false

	a.Router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000" || origin == ""
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