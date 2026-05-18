package handler

import (
	"github.com/bellapacx/kids-utopia/internal/users/dto"
	"github.com/bellapacx/kids-utopia/internal/users/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}
func (h *UserHandler) Me(c *gin.Context) {

	userID := c.GetString("userID")

	user, err := h.service.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": user})
}
func (h *UserHandler) UpdateMe(c *gin.Context) {

	userID := c.GetString("userID")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.UpdateMe(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}