package watermillzap

import (
	"github.com/ThreeDotsLabs/watermill"
	"go.uber.org/zap"
)

// Logger implements watermill.LoggerAdapter.
type Logger struct {
	zaplogger *zap.Logger
}

var _ watermill.LoggerAdapter = (*Logger)(nil)

func NewLogger(z *zap.Logger) Logger {
	return Logger{zaplogger: z}
}

func (l Logger) Error(msg string, err error, fields watermill.LogFields) {
	l.zaplogger.Error(msg, mapFields(fields)...)
}

func (l Logger) Info(msg string, fields watermill.LogFields) {
	l.zaplogger.Info(msg, mapFields(fields)...)
}

func (l Logger) Debug(msg string, fields watermill.LogFields) {
	l.zaplogger.Debug(msg, mapFields(fields)...)
}

func (l Logger) Trace(msg string, fields watermill.LogFields) {
	l.zaplogger.Debug(msg, mapFields(fields)...)
}

func (l Logger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &Logger{
		zaplogger: l.zaplogger.With(mapFields(fields)...),
	}
}

func mapFields(fields watermill.LogFields) []zap.Field {
	ret := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		ret = append(ret, zap.Any(k, v))
	}
	return ret
}
