package api

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	Level string `validate:"required,oneof=debug info warn error fatal"`
}

func NewLogger(config *LoggerConfig) (*zap.Logger, func(), error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig.LevelKey = "level"
	zapConfig.EncoderConfig.MessageKey = "message"
	zapConfig.Encoding = "json"

	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return nil, nil, err
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	logger, _ := zapConfig.Build()
	zap.ReplaceGlobals(logger)
	return logger, func() {
		_ = logger.Sync()
	}, nil
}
