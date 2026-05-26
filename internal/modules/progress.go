package modules

import (
	"github.com/gin-gonic/gin"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"

	progressHandler "github.com/bellapacx/kids-utopia/internal/progress/handler"
	progressRepo "github.com/bellapacx/kids-utopia/internal/progress/repository"
	progressRoutes "github.com/bellapacx/kids-utopia/internal/progress/routes"
	progressService "github.com/bellapacx/kids-utopia/internal/progress/service"

	"github.com/bellapacx/kids-utopia/pkg/database"
	"github.com/bellapacx/kids-utopia/pkg/middleware"
)

func RegisterProgress(
	r *gin.Engine,
	container *appcontainer.Container,
) {

	cfg := container.Config

	progRepo := progressRepo.NewProgressRepository(
		database.DB,
	)

	progService := progressService.NewProgressService(
		progRepo,
	)

	progHandler := progressHandler.NewProgressHandler(
		progService,
	)

	progressRoutes.RegisterRoutes(
		r.Group("/api/v1"),
		progHandler,
		middleware.AuthMiddleware(cfg.JWTSecret),
	)
}