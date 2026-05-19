package auth

import (
	"github.com/gin-gonic/gin"
)

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) *Routes {
	return &Routes{handler: handler}
}

func (r *Routes) Register(group *gin.RouterGroup) {

	auth := group.Group("/auth")
	{
		auth.POST("/register", r.handler.Register)
		auth.POST("/login", r.handler.Login)
		auth.POST("/verify-otp", r.handler.VerifyOTP)
		auth.POST("/refresh", r.handler.RefreshToken)
		auth.POST("/forgot-password", r.handler.ForgotPassword)
	auth.POST("/verify-reset-otp", r.handler.VerifyResetOTP)
	auth.POST("/reset-password", r.handler.ResetPassword)
	}
}