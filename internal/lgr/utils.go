package lgr

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey struct{}

var CtxKey = ctxKey{}

func GetLogger(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(CtxKey).(*zap.Logger); ok {
		return logger
	}
	return zap.L()
}
