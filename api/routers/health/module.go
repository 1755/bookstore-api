package health

import (
	"github.com/google/wire"
)

var Module = wire.NewSet(NewRouterBuilder, NewGetRouterBuilder)
