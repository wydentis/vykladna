package core_logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	*slog.Logger
	file *os.File
}

type LoggerContextKey struct{}

var (
	key = LoggerContextKey{}
)

func NewLogger(cfg Config) (*Logger, error) {
	if err := os.MkdirAll(cfg.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		cfg.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	writer := io.MultiWriter(os.Stdout, logFile)
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: cfg.Level,
	})

	return &Logger{
		Logger: slog.New(handler),
		file:   logFile,
	}, nil
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
		file:   l.file,
	}
}

func (l *Logger) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, l)
}

func FromContext(ctx context.Context) *Logger {
	logger, ok := ctx.Value(key).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return logger
}

func (l *Logger) Close() error {
	return l.file.Close()
}
