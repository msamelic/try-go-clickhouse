package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func New() *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.DisableStacktrace = true

	return zap.Must(config.Build())
}
