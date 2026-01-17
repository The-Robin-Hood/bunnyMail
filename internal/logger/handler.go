package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

// ColoredHandler with colors (for terminal)
type ColoredHandler struct {
	writer io.Writer
	level  slog.Level
	attrs  []slog.Attr
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[37m"
)

func NewColoredHandler(w io.Writer, level slog.Level) *ColoredHandler {
	return &ColoredHandler{
		writer: w,
		level:  level,
		attrs:  []slog.Attr{},
	}
}

func (h *ColoredHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *ColoredHandler) Handle(ctx context.Context, r slog.Record) error {
	// Time
	timeStr := r.Time.Format("2006-01-02 15:04:05.000")

	// Level with color
	var levelColor string
	var levelStr string

	switch r.Level {
	case slog.LevelDebug:
		levelColor = colorGray
		levelStr = "DEBUG"
	case slog.LevelInfo:
		levelColor = colorBlue
		levelStr = "INFOS"
	case slog.LevelWarn:
		levelColor = colorYellow
		levelStr = "ALERT"
	case slog.LevelError:
		levelColor = colorRed
		levelStr = "ERROR"
	default:
		levelColor = colorReset
		levelStr = r.Level.String()
	}

	var buf strings.Builder

	// [timestamp] [colored level] message
	fmt.Fprintf(&buf, "[%s]  %s[%s]%s  %s",
		timeStr,
		levelColor,
		levelStr,
		colorReset,
		r.Message,
	)

	// Add attributes
	r.Attrs(func(a slog.Attr) bool {
		buf.WriteString(fmt.Sprintf("  %s%s%s=%v",
			colorGray, a.Key, colorReset, a.Value))
		return true
	})

	for _, attr := range h.attrs {
		buf.WriteString(fmt.Sprintf("  %s%s%s=%v",
			colorGray, attr.Key, colorReset, attr.Value))
	}

	buf.WriteString("\n")

	_, err := h.writer.Write([]byte(buf.String()))
	return err
}

func (h *ColoredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &ColoredHandler{
		writer: h.writer,
		level:  h.level,
		attrs:  newAttrs,
	}
}

func (h *ColoredHandler) WithGroup(name string) slog.Handler {
	return h
}
