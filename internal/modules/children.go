package modules

import (
	"github.com/gin-gonic/gin"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"

	childrenRoutes "github.com/bellapacx/kids-utopia/internal/children"
	childrenHandler "github.com/bellapacx/kids-utopia/internal/children/handler"
	childrenRepo "github.com/bellapacx/kids-utopia/internal/children/repository"
	childrenService "github.com/bellapacx/kids-utopia/internal/children/service"

	"github.com/bellapacx/kids-utopia/pkg/database"
	"github.com/bellapacx/kids-utopia/pkg/middleware"
)

func RegisterChildren(
	r *gin.Engine,
	container *appcontainer.Container,
) {

	cfg := container.Config

	childRepo := childrenRepo.NewChildRepository(
		database.DB,
	)

	childService := childrenService.NewChildService(
		childRepo,
	)

	childHandler := childrenHandler.NewChildHandler(
		childService,
	)

	childrenRoutes.RegisterRoutes(
		r.Group("/api/v1"),
		childHandler,
		middleware.AuthMiddleware(cfg.JWTSecret),
	)
}