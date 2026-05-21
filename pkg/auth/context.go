package auth

import (
	"github.com/bellapacx/kids-utopia/pkg/contextkeys"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) string {
	v, ok := c.Get(contextkeys.UserID)
	if !ok || v == nil {
		return ""
	}
	return v.(string)
}

func GetRole(c *gin.Context) string {
	v, ok := c.Get(contextkeys.Role)
	if !ok || v == nil {
		return ""
	}
	return v.(string)
}