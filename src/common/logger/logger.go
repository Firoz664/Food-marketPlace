package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func init() {
	// Define custom colors
	yellow := "\033[33m"
	reset := "\033[0m"

	// Define encoder configuration
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(yellow + l.CapitalString() + reset)
	}

	// Initialize logger
	var err error
	Log, err = config.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}
