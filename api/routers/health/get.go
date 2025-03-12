package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetRouterBuilder struct {
}

func NewGetRouterBuilder() *GetRouterBuilder {
	return &GetRouterBuilder{}
}

func (r *GetRouterBuilder) Build(g *gin.RouterGroup) {
	g.GET("/", r.getHandler())
}

func (r *GetRouterBuilder) getHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	}
}
