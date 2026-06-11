package modules

import (
	"github.com/gin-gonic/gin"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"
	"github.com/bellapacx/kids-utopia/pkg/middleware"

	// SUBSCRIPTIONS
	subrepo "github.com/bellapacx/kids-utopia/internal/subscriptions/repository"
	subservice "github.com/bellapacx/kids-utopia/internal/subscriptions/service"

	accesssvc "github.com/bellapacx/kids-utopia/internal/access/service"

	bookmarkhandler "github.com/bellapacx/kids-utopia/internal/bookmarks/handler"
	bookmarkrepo "github.com/bellapacx/kids-utopia/internal/bookmarks/repository"
	bookmarkroutes "github.com/bellapacx/kids-utopia/internal/bookmarks/routes"
	bookmarkservice "github.com/bellapacx/kids-utopia/internal/bookmarks/service"
)

func RegisterBookmarks(
	r *gin.Engine,
	container *appcontainer.Container,
	
) {
	db := container.DB
    	cfg := container.Config

	// =========================
	// repositories & services
	// =========================
	subRepo := subrepo.New(db)
	subService := subservice.New(subRepo)

	// access layer
	accessSvc := accesssvc.New(subService)
	

	bookmarkRepo := bookmarkrepo.New(db)
	bookmarkService := bookmarkservice.New(bookmarkRepo, accessSvc)
	bookmarkHandler := bookmarkhandler.New(bookmarkService)

	bookmarkRoutes := bookmarkroutes.NewRoutes(bookmarkHandler)

	// =========================
	// 🔐 AUTH GROUP (IMPORTANT FIX)
	// =========================
	api := r.Group("/api/v1")

	// DO NOT forget this middleware
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret),) // or middleware.Auth(cfg.JWTSecret)

	bookmarkRoutes.Register(api)
}