package book

import "github.com/google/wire"

var Module = wire.NewSet(
	NewBasicDAO,
	NewBasicService,
	wire.Bind(new(DAO), new(*BasicDAO)),
	wire.Bind(new(Service), new(*BasicService)),
)
