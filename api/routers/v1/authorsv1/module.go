package authorsv1

import (
	"github.com/google/wire"
)

var Module = wire.NewSet(NewRouterBuilder, NewListRouterBuilder, NewGetRouterBuilder, NewCreateRouterBuilder, NewDeleteRouterBuilder, NewUpdateRouterBuilder)
