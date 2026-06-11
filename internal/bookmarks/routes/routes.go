package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/bookmarks/handler"
)

type Routes struct {
	handler *handler.Handler
}

func NewRoutes(h *handler.Handler) *Routes {
	return &Routes{handler: h}
}

func (r *Routes) Register(api *gin.RouterGroup) {
	bookmarks := api.Group("/bookmarks")

	bookmarks.POST("/create", r.handler.Create)
	bookmarks.DELETE("/delete", r.handler.Delete)

	bookmarks.GET("/book", r.handler.ListByBook)
	bookmarks.GET("/child", r.handler.ListByChild)
	bookmarks.GET("/child/detailed", r.handler.ListDetailedByChild)
}