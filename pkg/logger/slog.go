package logger

import (
	"errors"
	"log/slog"
	"os"
	"strings"
)

// slogLogger — адаптер slog.Logger, реализующий интерфейс Logger.
type slogLogger struct {
	logger *slog.Logger
}

// NewLogger создает новый Logger с указанным уровнем логирования (debug, info, warn, error).
func NewLogger(level string) (Logger, error) {
	lvl, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	return &slogLogger{
		logger: slog.New(handler),
	}, nil
}

func (l *slogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *slogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *slogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *slogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *slogLogger) Fatal(msg string, args ...any) {
	l.logger.Error(msg, args...)
	os.Exit(1)
}

func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{
		logger: l.logger.With(args...),
	}
}

func (l *slogLogger) WithOp(op string) Logger {
	return l.With(LogKeyOp, op)
}

func (l *slogLogger) WithRequestID(id string) Logger {
	return l.With(LogKeyRequestID, id)
}

func (l *slogLogger) WithUserID(id int64) Logger {
	return l.With(LogKeyUserID, id)
}

func (l *slogLogger) WithError(err error) Logger {
	return l.With(LogKeyError, err)
}

// parseLevel converts a string into a slog.Level.
func parseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, errors.New("invalid log level: " + s)
	}
}
