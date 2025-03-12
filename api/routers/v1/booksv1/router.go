package booksv1

import (
	"github.com/gin-gonic/gin"
)

type RouterBuilder struct {
	listRouter   *ListRouterBuilder
	getRouter    *GetRouterBuilder
	createRouter *CreateRouterBuilder
	deleteRouter *DeleteRouterBuilder
	updateRouter *UpdateRouterBuilder
}

func NewRouterBuilder(listRouter *ListRouterBuilder, getRouter *GetRouterBuilder, createRouter *CreateRouterBuilder, deleteRouter *DeleteRouterBuilder, updateRouter *UpdateRouterBuilder) *RouterBuilder {
	return &RouterBuilder{listRouter, getRouter, createRouter, deleteRouter, updateRouter}
}

func (r *RouterBuilder) Build(g *gin.RouterGroup) {
	root := g.Group("/v1/books")

	r.listRouter.Build(root)
	r.getRouter.Build(root)
	r.createRouter.Build(root)
	r.deleteRouter.Build(root)
	r.updateRouter.Build(root)
}
