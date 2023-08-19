package logging

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerKey = struct{}{}

// InitializeLogger returns a logger instance.
func InitializeLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "", // 時刻情報を表示しない
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     nil, // 時刻情報のエンコードを省略
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // ファイル名と行数を表示
	}

	atom := zap.NewAtomicLevelAt(zap.InfoLevel)

	config := zap.Config{
		Level:            atom,
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	l, err := config.Build(zap.AddCaller()) // AddCallerオプションでcaller情報を追加
	if err != nil {
		panic(err)
	}
	return l
}

// ContextWithLogger returns a new context with the logger.
func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// LoggerFromContext returns a logger from the context.
func LoggerFromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		fmt.Println("logger not found in context")
		return zap.L() // デフォルトロガーを返す
	}
	return logger
}
