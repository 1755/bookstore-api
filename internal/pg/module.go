package pg

import "github.com/google/wire"

var Module = wire.NewSet(
	NewBasicPool,
)
