package modules

import (
	"context"

	"github.com/gin-gonic/gin"

	// BOOKS
	bookrepo "github.com/bellapacx/kids-utopia/internal/books/repository"
	bookservice "github.com/bellapacx/kids-utopia/internal/books/service"

	// ACCESS
	accessservice "github.com/bellapacx/kids-utopia/internal/access/service"

	// SUBSCRIPTIONS
	subrepo "github.com/bellapacx/kids-utopia/internal/subscriptions/repository"
	subservice "github.com/bellapacx/kids-utopia/internal/subscriptions/service"

	// PROGRESS
	progressrepo "github.com/bellapacx/kids-utopia/internal/progress/repository"
	progressservice "github.com/bellapacx/kids-utopia/internal/progress/service"

	// READER ENGINE
	readerengine "github.com/bellapacx/kids-utopia/internal/reader/engine"
	"github.com/bellapacx/kids-utopia/internal/reader/events"
	readerroutes "github.com/bellapacx/kids-utopia/internal/reader/routes"

	// READER LIBRARY
	readerlibrary "github.com/bellapacx/kids-utopia/internal/reader/library"
	readerlibraryrepo "github.com/bellapacx/kids-utopia/internal/reader/library/repository"

	streakrepo "github.com/bellapacx/kids-utopia/internal/reader/streak/repository"
	streakservice "github.com/bellapacx/kids-utopia/internal/reader/streak/service"

	// READER SESSION
	sessionrepo "github.com/bellapacx/kids-utopia/internal/reader_session/repository"
	sessionservice "github.com/bellapacx/kids-utopia/internal/reader_session/service"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"

	"github.com/bellapacx/kids-utopia/pkg/middleware"
)

func RegisterReader(
	r *gin.Engine,
	container *appcontainer.Container,
) {

	cfg := container.Config

	// =========================
	// SUBSCRIPTIONS
	// =========================

	subRepo := subrepo.New(container.DB)

	subService := subservice.New(subRepo)

	// =========================
	// ACCESS
	// =========================

	accessService := accessservice.New(subService)

	// =========================
	// BOOKS
	// =========================

	bookRepo := bookrepo.NewBookRepository()

	bookService := bookservice.NewBookService(
		bookRepo,
		container.Storage,
		container.Queue,
		accessService,
	)

	// =========================
	// PROGRESS
	// =========================

	progressRepo := progressrepo.NewProgressRepository(container.DB)

	progressService := progressservice.NewProgressService(progressRepo)

	// =========================
	// SESSION
	// =========================

	sessionRepo := sessionrepo.New(container.DB)

	sessionService := sessionservice.New(sessionRepo)

	// =========================
	// STREAK (NEW)
	// =========================

	streakRepo := streakrepo.New(container.DB)

	streakService := streakservice.New(streakRepo)

	// =========================
	// EVENT BUS (NEW)
	// =========================

	eventBus := events.NewBus()

	// =========================
	// ENGINE
	// =========================

	engine := readerengine.New(
		accessService,
		bookService,
		sessionService,
		progressService,
		streakService, // ✅ FIX 1
		eventBus,      // ✅ FIX 2
	)

	// =========================
	// EVENT SUBSCRIPTIONS
	// =========================

	eventBus.Subscribe(events.ProgressUpdated, func(e events.Event) {
		_ = streakService.UpdateStreak(context.Background(), e.ChildID)
	})

	// =========================
	// LIBRARY
	// =========================

	libraryRepo := readerlibraryrepo.New(container.DB)

	libraryService := readerlibrary.New(
		libraryRepo,
		bookService,
	)

	libraryHandler := readerlibrary.NewHandler(libraryService)

	// =========================
	// ROUTES
	// =========================

	api := r.Group("/api/v1")

	api.Use(
		middleware.AuthMiddleware(cfg.JWTSecret),
	)

	readerroutes.RegisterRoutes(api, engine)

	readerlibrary.RegisterRoutes(api, libraryHandler)
}