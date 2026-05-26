package modules

import (
	"github.com/gin-gonic/gin"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"

	usersRoutes "github.com/bellapacx/kids-utopia/internal/users"
	usersHandler "github.com/bellapacx/kids-utopia/internal/users/handler"
	usersRepo "github.com/bellapacx/kids-utopia/internal/users/repository"
	usersService "github.com/bellapacx/kids-utopia/internal/users/service"

	"github.com/bellapacx/kids-utopia/pkg/database"
	"github.com/bellapacx/kids-utopia/pkg/middleware"
)

func RegisterUsers(
	r *gin.Engine,
	container *appcontainer.Container,
) {

	cfg := container.Config

	userRepo := usersRepo.NewUserRepository(
		database.DB,
	)

	userService := usersService.NewUserService(
		userRepo,
	)

	userHandler := usersHandler.NewUserHandler(
		userService,
	)

	usersRoutes.RegisterRoutes(
		r.Group("/api/v1"),
		userHandler,
		middleware.AuthMiddleware(cfg.JWTSecret),
	)
}