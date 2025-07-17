package utils

import "go.uber.org/zap"

func SafeLog(logger *zap.Logger, level string, msg string, fields ...zap.Field) {
	if logger == nil {
		return
	}
	switch level {
	case "info":
		logger.Info(msg, fields...)
	case "error":
		logger.Error(msg, fields...)
	case "warn":
		logger.Warn(msg, fields...)
	default:
		logger.Info(msg, fields...)
	}
}
