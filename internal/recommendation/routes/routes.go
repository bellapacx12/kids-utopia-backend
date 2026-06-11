package routes

import (
	"github.com/bellapacx/kids-utopia/internal/recommendation/handler"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	handler *handler.Handler
}

func NewRoutes(h *handler.Handler) *Routes {
	return &Routes{handler: h}
}

func (r *Routes) Register(api *gin.RouterGroup) {

	rec := api.Group("/recommendations")

	rec.GET("", r.handler.Recommend)
}