package modules

import (
	"github.com/gin-gonic/gin"

	appcontainer "github.com/bellapacx/kids-utopia/internal/app"

	"github.com/bellapacx/kids-utopia/internal/auth"
	"github.com/bellapacx/kids-utopia/internal/notifications/email"
	"github.com/bellapacx/kids-utopia/internal/notifications/otp"
	"github.com/bellapacx/kids-utopia/internal/notifications/sms"
)

func RegisterAuth(
	r *gin.Engine,
	container *appcontainer.Container,
) {

	cfg := container.Config

	emailSender := email.NewSendGrid(
		cfg.SendGridAPIKey,
		cfg.FromEmail,
	)

	smsSender := sms.NewSender()

	otpRouter := otp.NewRouter(
		emailSender,
		smsSender,
	)

	otpService := otp.NewService(
		otpRouter,
	)

	authRepo := &auth.Repository{}

	authService := auth.NewService(
		authRepo,
		otpService,
		cfg.JWTSecret,
	)

	authHandler := auth.NewHandler(
		authService,
	)

	auth.NewRoutes(authHandler).
		Register(r.Group("/api/v1"))
}