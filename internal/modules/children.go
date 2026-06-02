package modules

import (
	"github.com/gin-gonic/gin"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"

	childrenRoutes "github.com/bellapacx/kids-utopia/internal/children"
	childrenHandler "github.com/bellapacx/kids-utopia/internal/children/handler"
	childrenRepo "github.com/bellapacx/kids-utopia/internal/children/repository"
	childrenService "github.com/bellapacx/kids-utopia/internal/children/service"

	streakhandler "github.com/bellapacx/kids-utopia/internal/streak/handler"
	streakrepo "github.com/bellapacx/kids-utopia/internal/streak/repository"
	streakservice "github.com/bellapacx/kids-utopia/internal/streak/service"

	analyticshandler "github.com/bellapacx/kids-utopia/internal/analytics/handler"
	analyticsrepo "github.com/bellapacx/kids-utopia/internal/analytics/repository"
	analyticsservice "github.com/bellapacx/kids-utopia/internal/analytics/service"

	// READER SESSION
	sessionrepo "github.com/bellapacx/kids-utopia/internal/reader_session/repository"

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

	sessionRepo := sessionrepo.New(container.DB)
	

streakRepo := streakrepo.New(container.DB)
streakService := streakservice.New(streakRepo)

streakHandler := streakhandler.New(streakService) // ✅ correct

	// =========================
	// EVENT BUS (NEW)
	// =========================
analyticsRepo := analyticsrepo.New(container.DB)
analyticsService := analyticsservice.New(analyticsRepo,streakRepo, sessionRepo)

analyticsHandler := analyticshandler.New(analyticsService) // ✅ correc

	childrenRoutes.RegisterRoutes(
		r.Group("/api/v1"),
		childHandler,
		middleware.AuthMiddleware(cfg.JWTSecret),
		streakHandler,
		analyticsHandler,
	)
}