package logger

import (
	"fmt"

	"github.com/goletan/security/shared/scrubber"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger   *zap.Logger
	scrubber *scrubber.Scrubber
}

func NewLogger() (*ZapLogger, error) {
	cfg := zap.NewProductionConfig()

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	scrub := scrubber.NewScrubber()

	return &ZapLogger{
		logger:   logger,
		scrubber: scrub,
	}, nil
}

func (z *ZapLogger) Info(message string, fields ...zap.Field) {
	z.scrubMessage(message, fields...)
	z.logger.Info(message, fields...)
}

func (z *ZapLogger) Debug(message string, fields ...zap.Field) {
	z.scrubMessage(message, fields...)
	z.logger.Debug(message, fields...)
}

func (z *ZapLogger) Warn(message string, fields ...zap.Field) {
	z.scrubMessage(message, fields...)
	z.logger.Warn(message, fields...)
}

func (z *ZapLogger) Error(message string, fields ...zap.Field) {
	z.scrubMessage(message, fields...)
	z.logger.Error(message, fields...)
}

func (z *ZapLogger) Fatal(message string, fields ...zap.Field) {
	z.scrubMessage(message, fields...)
	z.logger.Fatal(message, fields...)
}

func (z *ZapLogger) WithContext(context map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(context))
	for k, v := range context {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

// scrubMessage scrubs the message and fields to remove sensitive data.
func (z *ZapLogger) scrubMessage(message string, fields ...zap.Field) {
	message = z.scrubber.Scrub(message)

	fieldCopy := make([]zap.Field, len(fields))
	copy(fieldCopy, fields)

	for i := range fieldCopy {
		if fieldCopy[i].Key != "" {
			originalValue := fieldCopy[i].String
			switch fieldCopy[i].Type {
			case zapcore.StringType:
				fieldCopy[i] = zap.String(fieldCopy[i].Key, z.scrubber.Scrub(originalValue))
			case zapcore.ErrorType:
				if err, ok := fieldCopy[i].Interface.(error); ok {
					fieldCopy[i] = zap.String(fieldCopy[i].Key, z.scrubber.Scrub(err.Error()))
				}
			case zapcore.ReflectType, zapcore.ObjectMarshalerType, zapcore.StringerType:
				fields[i] = zap.String(fields[i].Key, z.scrubber.Scrub(fmt.Sprintf("%v", fields[i].Interface)))
			default:
				fields[i] = zap.Any(fields[i].Key, z.scrubber.Scrub(fmt.Sprintf("%v", fields[i].Interface)))
			}
		}
	}
}
