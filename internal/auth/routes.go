package auth

import (
	"github.com/bellapacx/kids-utopia/pkg/config"
	"github.com/bellapacx/kids-utopia/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) *Routes {
	return &Routes{handler: handler}
}

func (r *Routes) Register(group *gin.RouterGroup) {


	cfg := config.Load()

	auth := group.Group("/auth")
	{
		auth.POST("/register", r.handler.Register)
		auth.POST("/login", r.handler.Login)
		auth.POST("/verify-otp", r.handler.VerifyOTP)
		auth.POST("/refresh", r.handler.RefreshToken)
		auth.POST("/forgot-password", r.handler.ForgotPassword)
	auth.POST("/reset-password", r.handler.ResetPassword)
	auth.POST("/logout", r.handler.Logout)
	auth.POST("/resend-otp", r.handler.ResendOTP)
	auth.POST("/verify-phone", r.handler.VerifyPhone)
	auth.POST("/verify-email", r.handler.VerifyEmail)
	auth.POST("/send-email-otp", r.handler.SendEmailOTP)
	auth.POST("/send-phone-otp", r.handler.SendPhoneOTP)
	}
	// PROTECTED ROUTES
	protected := group.Group("/auth")

	protected.Use(
		middleware.AuthMiddleware(cfg.JWTSecret),
	)

	{
		protected.GET(
			"/verification-session",
			r.handler.VerificationSession,
		)
	}

}