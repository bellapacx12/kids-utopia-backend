package engine

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/pkg/contextkeys"
)

func (e *Engine) OpenHandler(c *gin.Context) {

	var req OpenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.GetString(contextkeys.UserID)

	state, err := e.Open(
		c.Request.Context(),
		userID,
		req.ChildID,
		req.BookID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, state)
}

func (e *Engine) UpdateHandler(c *gin.Context) {

	var req UpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.GetString(contextkeys.UserID)

	err := e.Update(
		c.Request.Context(),
		userID,
		req.ChildID,
		req.BookID,
		req.Page,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "updated",
	})
}

func (e *Engine) CloseHandler(c *gin.Context) {

	var req CloseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.GetString(contextkeys.UserID)

	err := e.Close(
		c.Request.Context(),
		userID,
		req.ChildID,
		req.BookID,
		req.Page,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "closed",
	})
}
func (e *Engine) StateHandler(c *gin.Context) {

	bookID := c.Param("bookId")
	childID := c.Param("childId")

	userID := c.GetString(contextkeys.UserID)

	state, err := e.State(
		c.Request.Context(),
		userID,
		childID,
		bookID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, state)
}