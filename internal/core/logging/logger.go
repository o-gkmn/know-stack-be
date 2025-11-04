package logging

import (
	"context"
	"fmt"
	"knowstack/internal/core/config"
	"log/slog"
	"os"
	"strings"
	"time"
)

type CustomHandler struct {
	slog.Handler
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	timeStr := r.Time.Format("15:04:05")

	levelColor := "\033[32m"
	level := "INFO"
	switch r.Level {
	case slog.LevelError:
		levelColor = "\033[31m"
		level = "ERROR"
	case slog.LevelWarn:
		levelColor = "\033[33m"
		level = "WARN"
	case slog.LevelDebug:
		levelColor = "\033[36m"
		level = "DEBUG"
	}

	fmt.Printf("%s %s[%s]\033[0m %s", timeStr, levelColor, level, r.Message)

	r.Attrs(func(a slog.Attr) bool {
		// Colorize keys (magenta) while leaving values default
		fmt.Printf("  \033[35m%s\033[0m = %v", strings.ToUpper(a.Key), a.Value.Any())
		return true
	})

	fmt.Println()

	return nil
}

func Init(cfg config.Logger) *slog.Logger {
	level := parseLogLevel(cfg.Level)

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	}

	h := chooseHandler(cfg.Format, opts)
	l := slog.New(h)

	slog.SetDefault(l)

	return l
}

func chooseHandler(format string, opts *slog.HandlerOptions) slog.Handler {
	switch strings.ToLower(format) {
	case "json":
		return slog.NewJSONHandler(os.Stdout, opts)
	case "custom":
		return &CustomHandler{
			Handler: slog.NewTextHandler(os.Stdout, opts),
		}
	default:
		opts.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(time.Now().Format(time.RFC3339))
			}

			return a
		}

		return slog.NewTextHandler(os.Stdout, opts)
	}
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
