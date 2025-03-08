//go:build wireinject
// +build wireinject

package api

import (
	"github.com/1755/bookstore-api/api/routers/health"
	"github.com/google/wire"
)

func InjectApplication(
	configPath ConfigPath,
) (*Application, func(), error) {
	panic(wire.Build(
		NewApplication,
		NewServer,
		ConfigModule,
		health.Module,
	))
}
