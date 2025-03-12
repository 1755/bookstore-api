//go:build tools

package main

// This prevents `go mod tidy` from purging required tools.
import (
	_ "github.com/swaggo/swag/cmd/swag"
)
