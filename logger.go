package gologger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger offers a simple to use logger interface.
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Err(err error)

	WithField(name string, value any) Logger
	WithFields(map[string]any) Logger

	Flush() error
}

// LogLevel defines the log level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelError
)

type logger struct {
	zap    *zap.Logger
	fields map[string]any
}

// New assembles a new logger with the given log level.
// "skipCallers" defines how many callers are skipped in the hierarchy for displaying the caller.
func New(loglevel LogLevel, skipCallers int) Logger {
	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	zapLogger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	), zap.AddCallerSkip(1+skipCallers))
	atom.SetLevel(logLevelToZap(loglevel))

	return &logger{
		zap:    zapLogger,
		fields: make(map[string]any),
	}
}

func logLevelToZap(loglevel LogLevel) zapcore.Level {
	switch loglevel {
	case LogLevelDebug:
		return zap.DebugLevel
	case LogLevelInfo:
		return zap.InfoLevel
	case LogLevelError:
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func (l *logger) Flush() error {
	return l.zap.Sync()
}

func (l *logger) Debugf(format string, args ...any) {
	l.zap.With(map2fields(l.fields)...).Debug(stringify(format, args...))
}

func (l *logger) Infof(format string, args ...any) {
	l.zap.With(map2fields(l.fields)...).Info(stringify(format, args...))
}

func (l *logger) Errorf(format string, args ...any) {
	l.zap.With(map2fields(l.fields)...).Error(stringify(format, args...))
}

func (l *logger) Err(err error) {
	l.zap.With(map2fields(l.fields)...).Error(err.Error())
}

func (l *logger) WithField(name string, value any) Logger {
	newLogger := l.copy()
	newLogger.fields[name] = value

	return newLogger
}

func (l *logger) WithFields(fields map[string]any) Logger {
	newLogger := l.copy()
	for key, value := range fields {
		newLogger.fields[key] = value
	}

	return newLogger
}

func (l *logger) copy() *logger {
	newLogger := logger{
		zap:    l.zap,
		fields: make(map[string]any),
	}

	for key, value := range l.fields {
		newLogger.fields[key] = value
	}

	return &newLogger
}

func map2fields(m map[string]any) []zap.Field {
	fields := make([]zap.Field, 0)
	for k, v := range m {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

func stringify(format string, args ...any) string {
	if len(args) == 0 {
		return format
	}

	return fmt.Sprintf(format, args...)
}
