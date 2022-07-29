package logging

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func MustNewLogger(minLevel string) *zap.SugaredLogger {
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapLevel(minLevel))
	zc.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	logger, err := zc.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}

func zapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}
