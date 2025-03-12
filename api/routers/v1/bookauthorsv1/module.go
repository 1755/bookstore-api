package bookauthorsv1

import (
	"github.com/google/wire"
)

var Module = wire.NewSet(NewRouterBuilder, NewListRouterBuilder, NewCreateRouterBuilder, NewDeleteRouterBuilder)
