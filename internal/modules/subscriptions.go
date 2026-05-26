package modules

import (
	"github.com/gin-gonic/gin"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"

	subhandler "github.com/bellapacx/kids-utopia/internal/subscriptions/handler"
	subrepo "github.com/bellapacx/kids-utopia/internal/subscriptions/repository"
	subroutes "github.com/bellapacx/kids-utopia/internal/subscriptions/routes"
	subservice "github.com/bellapacx/kids-utopia/internal/subscriptions/service"

	"github.com/bellapacx/kids-utopia/pkg/database"
)

func RegisterSubscriptions(
	r *gin.Engine,
	container *appcontainer.Container,
) {

	subRepo := subrepo.New(
		database.DB,
	)

	subService := subservice.New(
		subRepo,
	)

	subHandler := subhandler.New(
		subService,
	)

	subroutes.RegisterSubscriptionRoutes(
		r.Group("/api/v1/subscriptions"),
		subHandler,
	)
}