//go:build tools

package main

// This prevents `go mod tidy` from purging required tools.
import (
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/swaggo/swag/cmd/swag"
	_ "github.com/vektra/mockery/v2"
)
