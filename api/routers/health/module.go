package health

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var Module = wire.NewSet(NewRouterBuilder)

type RouterBuilder struct {
}

func NewRouterBuilder() *RouterBuilder {
	return &RouterBuilder{}
}

func (r *RouterBuilder) Build(g *gin.RouterGroup) {
	g.GET("/health", r.getHandler())
}
