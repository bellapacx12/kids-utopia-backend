package modules

import (
	"github.com/gin-gonic/gin"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"

	// BOOKS
	bookrepo "github.com/bellapacx/kids-utopia/internal/books/repository"

	// PROGRESS
	progressrepo "github.com/bellapacx/kids-utopia/internal/progress/repository"

	// BOOKMARKS
	bookmarkrepo "github.com/bellapacx/kids-utopia/internal/bookmarks/repository"
	childrepo "github.com/bellapacx/kids-utopia/internal/children/repository"

	// RECOMMENDATION
	recommendationhandler "github.com/bellapacx/kids-utopia/internal/recommendation/handler"
	recommendationroutes "github.com/bellapacx/kids-utopia/internal/recommendation/routes"
	recommendationservice "github.com/bellapacx/kids-utopia/internal/recommendation/service"

	// MIDDLEWARE
	"github.com/bellapacx/kids-utopia/pkg/middleware"
)
func RegisterRecommendation(
	r *gin.Engine,
	container *appcontainer.Container,
) {

	// =========================
	// repositories
	// =========================
	bookRepo := bookrepo.NewBookRepository()
	progressRepo := progressrepo.NewProgressRepository(container.DB)
	bookmarkRepo := bookmarkrepo.New(container.DB)
	childRepo := childrepo.NewChildRepository(container.DB)

	// =========================
	// service
	// =========================
	service := recommendationservice.New(
		bookRepo,
		progressRepo,
		bookmarkRepo,
		childRepo,
	)

	// =========================
	// handler
	// =========================
	handler := recommendationhandler.New(service)

	// =========================
	// routes
	// =========================
	routes := recommendationroutes.NewRoutes(handler)

	// =========================
	// auth protected group
	// =========================
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(container.Config.JWTSecret))

	routes.Register(api)
}