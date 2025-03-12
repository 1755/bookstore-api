//go:build wireinject
// +build wireinject

package api

import (
	"github.com/1755/bookstore-api/api/routers/health"
	"github.com/1755/bookstore-api/api/routers/v1/authorsv1"
	"github.com/1755/bookstore-api/api/routers/v1/bookauthorsv1"
	"github.com/1755/bookstore-api/api/routers/v1/booksv1"
	"github.com/1755/bookstore-api/internal/author"
	"github.com/1755/bookstore-api/internal/book"
	"github.com/1755/bookstore-api/internal/pg"
	"github.com/google/wire"
)

func InjectApplication(
	configPath ConfigPath,
) (*Application, func(), error) {
	panic(wire.Build(
		NewLogger,
		NewContext,
		NewMonitoring,
		NewApplication,
		NewServer,
		ConfigModule,
		health.Module,
		booksv1.Module,
		authorsv1.Module,
		bookauthorsv1.Module,
		book.Module,
		author.Module,
		pg.Module,
	))
}
