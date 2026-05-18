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

	err := h.service.Register(req)
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

	res, err := h.service.Login(req, deviceID)
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

	err := h.service.VerifyOTP(req)
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

	res, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.Success(c, res)
}