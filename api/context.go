package api

import (
	"context"

	"github.com/1755/bookstore-api/internal/lgr"
	"go.uber.org/zap"
)

func NewContext(logger *zap.Logger) context.Context {
	return context.WithValue(context.Background(), lgr.CtxKey, logger)
}
