package modules

import (
	"github.com/gin-gonic/gin"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"

	readerhandler "github.com/bellapacx/kids-utopia/internal/reader_session/handler"
	readerrepo "github.com/bellapacx/kids-utopia/internal/reader_session/repository"
	readerroutes "github.com/bellapacx/kids-utopia/internal/reader_session/routes"
	readerservice "github.com/bellapacx/kids-utopia/internal/reader_session/service"

	"github.com/bellapacx/kids-utopia/pkg/middleware"
)

func RegisterReaderSession(
	r *gin.Engine,
	container *appcontainer.Container,
) {

	cfg := container.Config

	repo := readerrepo.New(
		container.DB,
	)

	service := readerservice.New(
		repo,
	)

	handler := readerhandler.New(
		service,
	)

	api := r.Group("/api/v1")

	api.Use(
		middleware.AuthMiddleware(cfg.JWTSecret),
	)

	readerroutes.RegisterRoutes(
		api,
		handler,
	)
}