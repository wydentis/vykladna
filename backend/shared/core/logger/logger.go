package core_logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
	file *os.File
}

func NewLogger(cfg Config) (*Logger, error) {
	file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	writer := io.MultiWriter(os.Stdout, file)
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: cfg.Level,
	})

	return &Logger{
		Logger: slog.New(handler),
		file:   file,
	}, nil
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
		file:   l.file,
	}
}

func (l *Logger) Close() error {
	return l.file.Close()
}
