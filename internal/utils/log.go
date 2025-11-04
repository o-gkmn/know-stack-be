package utils

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func LogInfo(msg string, args ...any) {
	slog.Info(msg, args...)
}

func LogDebug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func LogWarn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func LogError(msg string, args ...any) {
	slog.Error(msg, args...)
}

func LogFatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func LogFatalWithErr(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, "error", err)
	}
	LogFatal(msg, args...)
}

func LogErrorWithErr(msg string, err error, args ...any) {
	if err != nil {
		args = append(args, "error", err)
	}
	slog.Error(msg, args...)
}

func LogWithContext(ctx context.Context, level slog.Level, msg string, args ...any) {
	slog.Log(ctx, level, msg, args...)
}

func LogWithGinContext(ctx *gin.Context, level slog.Level, msg string, args ...any) {
	LogWithContext(ctx, level, msg, args...)
}

func LogHTTPRequest(method, path string, statusCode int, duration time.Duration, args ...any) {
	logArgs := []any{
		"method", method,
		"path", path,
		"status", statusCode,
		"duration", duration.String(),
	}
	logArgs = append(logArgs, args...)
	slog.Info("HTTP request", logArgs...)
}
