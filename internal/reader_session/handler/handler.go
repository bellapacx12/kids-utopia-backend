package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/reader_session/dto"
	"github.com/bellapacx/kids-utopia/internal/reader_session/service"
	"github.com/bellapacx/kids-utopia/pkg/contextkeys"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) StartSession(c *gin.Context) {

	var req dto.StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString(contextkeys.UserID)

	session, err := h.service.StartSession(
		c.Request.Context(),
		userID,
		req.ChildID,
		req.BookID,
		req.Page,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, session)
}

func (h *Handler) UpdateSession(c *gin.Context) {

	var req dto.UpdateSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateSession(
		c.Request.Context(),
		req.SessionID,
		req.Page,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "updated"})
}

func (h *Handler) EndSession(c *gin.Context) {

	var req dto.EndSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.EndSession(
		c.Request.Context(),
		req.SessionID,
		req.Page,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ended"})
}

func (h *Handler) GetSession(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing session id"})
		return
	}

	session, err := h.service.GetSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}
