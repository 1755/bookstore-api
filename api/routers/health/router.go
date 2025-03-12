package health

import (
	"github.com/gin-gonic/gin"
)

type RouterBuilder struct {
	getRouter *GetRouterBuilder
}

func NewRouterBuilder(getRouter *GetRouterBuilder) *RouterBuilder {
	return &RouterBuilder{getRouter}
}

func (r *RouterBuilder) Build(g *gin.RouterGroup) {
	root := g.Group("/health")

	r.getRouter.Build(root)
}
