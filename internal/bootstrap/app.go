package bootstrap

import (
	"github.com/gin-gonic/gin"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"
)

type App struct {
	Router    *gin.Engine
	Container *appcontainer.Container
}

func NewApp() *App {

	r := gin.Default()

	container := appcontainer.NewContainer()

	app := &App{
		Router:    r,
		Container: container,
	}

	app.registerMiddlewares()
	app.registerHealth()
	app.registerModules()

	return app
}