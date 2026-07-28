package core_logger_zap

import (
	"time"

	core_logger "github.com/wasstend/todoapp-golang/internal/core/logger"
	"go.uber.org/zap"
)

func toZapField(field core_logger.Field) zap.Field {
	switch value := field.Value.(type) {
	case string:
		return zap.String(field.Key, value)
	case int:
		return zap.Int(field.Key, value)
	case error:
		return zap.NamedError(field.Key, value)
	case time.Time:
		return zap.Time(field.Key, value)
	case time.Duration:
		return zap.Duration(field.Key, value)
	default:
		return zap.Any(field.Key, value)
	}
}

func toZapFields(fields []core_logger.Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))

	for _, field := range fields {
		zapFields = append(zapFields, toZapField(field))
	}

	return zapFields
}
