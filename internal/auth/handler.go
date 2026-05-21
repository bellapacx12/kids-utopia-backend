package auth

import (
	"github.com/bellapacx/kids-utopia/pkg/response"
	"github.com/bellapacx/kids-utopia/pkg/validator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}
func (h *Handler) Register(c *gin.Context) {

	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid request")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	err := h.service.Register(c.Request.Context(),req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "OTP sent",
	})
}
func (h *Handler) Login(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid request")
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	// ====================================
	// DEVICE ID (IMPORTANT FOR SESSIONS)
	// ====================================

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		deviceID = "unknown-device"
	}

	// ====================================
	// LOGIN SERVICE CALL
	// ====================================

	res, err := h.service.Login(c.Request.Context(),req, deviceID)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.Success(c, res)
}
func (h *Handler) VerifyOTP(c *gin.Context) {

	var req VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.VerifyOTP(c.Request.Context(),req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"verified": true,
		},
	})
}
func (h *Handler) RefreshToken(c *gin.Context) {

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid request")
		return
	}

	res, err := h.service.RefreshToken(c.Request.Context(),req.RefreshToken)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.Success(c, res)
}

func (h *Handler) ForgotPassword(c *gin.Context) {

	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.ForgotPassword(c.Request.Context(),req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "otp sent"})
}
func (h *Handler) VerifyResetOTP(c *gin.Context) {

	var req VerifyResetOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.VerifyResetOTP(c.Request.Context(),req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "otp verified"})
}
func (h *Handler) ResetPassword(c *gin.Context) {

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.ResetPassword(c.Request.Context(),req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "password updated"})
}
func (h *Handler) Logout(c *gin.Context) {

	var req LogoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.Logout(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "logged out successfully",
	})
}
func (h *Handler) ResendOTP(c *gin.Context) {

	var req ResendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.ResendOTP(c.Request.Context(), req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "otp resent",
	})
}
func (h *Handler) VerificationSession(c *gin.Context) {

	userID := c.GetString("user_id")

	res, err := h.service.VerificationSession(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, res)
}
func (h *Handler) VerifyEmail(c *gin.Context) {

	var req VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.VerifyEmail(
		c.Request.Context(),
		req,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "email verified successfully",
	})
}
func (h *Handler) VerifyPhone(c *gin.Context) {

	var req VerifyPhoneRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.VerifyPhone(
		c.Request.Context(),
		req,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "phone verified successfully",
	})
}