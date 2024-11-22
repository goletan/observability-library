package logger

import "go.uber.org/zap"

// LoggerInterface defines the interface for logging.
type LoggerInterface interface {
	Info(message string, fields ...zap.Field)
	Debug(message string, fields ...zap.Field)
	Warn(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	WithContext(context map[string]interface{}) []zap.Field
}

type ZapLogger struct {
	logger *zap.Logger
}

func NewLogger(cfg *zap.Config) (*ZapLogger, error) {
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{logger: logger}, nil
}

func (z *ZapLogger) Info(message string, fields ...zap.Field) {
	z.logger.Info(message, fields...)
}

func (z *ZapLogger) Debug(message string, fields ...zap.Field) {
	z.logger.Debug(message, fields...)
}

func (z *ZapLogger) Warn(message string, fields ...zap.Field) {
	z.logger.Warn(message, fields...)
}

func (z *ZapLogger) Error(message string, fields ...zap.Field) {
	z.logger.Error(message, fields...)
}

func (z *ZapLogger) Fatal(message string, fields ...zap.Field) {
	z.logger.Fatal(message, fields...)
}

func (z *ZapLogger) WithContext(context map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(context))
	for k, v := range context {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}
