package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps a zap.Logger for structured JSON logging
var Logger *zap.Logger

func InitLogger() error {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "datetime"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	Logger, err = cfg.Build()
	return err
}

func SyncLogger() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Usage:
//   logging.Logger.Info("message", zap.String("key", value))
//   logging.Logger.Error("error message", zap.Error(err))
