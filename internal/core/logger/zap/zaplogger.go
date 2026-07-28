package core_logger_zap

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	core_logger "github.com/wasstend/todoapp-golang/internal/core/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger

	file *os.File
}

func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()

	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("Unmarshal log level: %w\n", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("Mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	zap := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Logger{
		Logger: zap,
		file:   logFile,
	}, nil
}

func (l *Logger) Debug(msg string, fields ...core_logger.Field) {
	l.Logger.Debug(msg, toZapFields(fields)...)
}
func (l *Logger) Warn(msg string, fields ...core_logger.Field) {
	l.Logger.Warn(msg, toZapFields(fields)...)
}
func (l *Logger) Error(msg string, fields ...core_logger.Field) {
	l.Logger.Error(msg, toZapFields(fields)...)
}
func (l *Logger) Fatal(msg string, fields ...core_logger.Field) {
	l.Logger.Fatal(msg, toZapFields(fields)...)
}

func (l *Logger) With(fields ...core_logger.Field) core_logger.Logger {
	return &Logger{
		Logger: l.Logger.With(toZapFields(fields)...),
		file:   l.file,
	}
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close logger:", err)
	}
}
