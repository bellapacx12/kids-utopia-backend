package modules

import (
	"github.com/gin-gonic/gin"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"

	// ACCESS
	accessmiddleware "github.com/bellapacx/kids-utopia/internal/access/middleware"
	accessservice "github.com/bellapacx/kids-utopia/internal/access/service"

	// BOOKS
	bookhandler "github.com/bellapacx/kids-utopia/internal/books/handler"
	bookrepo "github.com/bellapacx/kids-utopia/internal/books/repository"
	bookroutes "github.com/bellapacx/kids-utopia/internal/books/routes"
	bookservice "github.com/bellapacx/kids-utopia/internal/books/service"

	// SUBSCRIPTIONS
	subrepo "github.com/bellapacx/kids-utopia/internal/subscriptions/repository"
	subservice "github.com/bellapacx/kids-utopia/internal/subscriptions/service"

	// INFRA
	"github.com/bellapacx/kids-utopia/pkg/database"
	"github.com/bellapacx/kids-utopia/pkg/middleware"
)

func RegisterBooks(
	r *gin.Engine,
	container *appcontainer.Container,
) {

	cfg := container.Config

	// =========================
	// SUBSCRIPTIONS
	// =========================
	subRepo := subrepo.New(
		database.DB,
	)

	subService := subservice.New(
		subRepo,
	)

	// =========================
	// ACCESS
	// =========================
	accessSvc := accessservice.New(
		subService,
	)

	// =========================
	// BOOKS
	// =========================
	bookRepo := bookrepo.NewBookRepository()

	bookService := bookservice.NewBookService(
		bookRepo,
		container.Storage,
		container.Queue,
		accessSvc,
	)

	bookHandler := bookhandler.NewBookHandler(
		bookService,
	)

	// =========================
	// ACCESS MIDDLEWARE
	// =========================
	accessMw := accessmiddleware.New(
		accessSvc,
		bookRepo,
	)

	// =========================
	// READER ROUTES
	// =========================
	readerGroup := r.Group("/api/v1/books")

	readerGroup.Use(
		middleware.AuthMiddleware(cfg.JWTSecret),
	)

	bookroutes.RegisterReaderRoutes(
		readerGroup,
		bookHandler,
		accessMw,
	)

	// =========================
	// EDITOR ROUTES
	// =========================
	editorBooks := r.Group("/api/v1/books")

	editorBooks.Use(
		middleware.AuthMiddleware(cfg.JWTSecret),
	)

	editorBooks.Use(
		middleware.RequireRoles("editor", "admin"),
	)

	bookroutes.RegisterEditorBookRoutes(
		editorBooks,
		bookHandler,
	)
}