package bookauthorsv1

import (
	"github.com/gin-gonic/gin"
)

type RouterBuilder struct {
	listRouter   *ListRouterBuilder
	createRouter *CreateRouterBuilder
	deleteRouter *DeleteRouterBuilder
}

func NewRouterBuilder(listRouter *ListRouterBuilder, createRouter *CreateRouterBuilder, deleteRouter *DeleteRouterBuilder) *RouterBuilder {
	return &RouterBuilder{listRouter, createRouter, deleteRouter}
}

func (r *RouterBuilder) Build(g *gin.RouterGroup) {
	root := g.Group("/v1/books/:id/authors")

	r.listRouter.Build(root)
	r.createRouter.Build(root)
	r.deleteRouter.Build(root)
}
