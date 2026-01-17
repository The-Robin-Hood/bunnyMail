package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var Log *slog.Logger

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type Config struct {
	Level      LogLevel
	OutputFile string
	UseJSON    bool // true = JSON format, false = text format
}

func Init(cfg Config) error {
	var level slog.Level
	switch cfg.Level {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo:
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var writer io.Writer = os.Stdout

	if cfg.OutputFile != "" {
		logDir := filepath.Dir(cfg.OutputFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		file, err := os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		writer = io.MultiWriter(os.Stdout, file)
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: level,
	}

	if cfg.UseJSON {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = NewColoredHandler(writer, level)
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)

	return nil
}

func With(args ...any) *slog.Logger {
	return Log.With(args...)
}

// WithContext returns a logger from context (for request tracking)
func WithContext(ctx context.Context) *slog.Logger {
	return Log
}

func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}
