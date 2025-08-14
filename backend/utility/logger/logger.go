package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type _logger struct {
	sugar *zap.SugaredLogger
}

var Logger *_logger

func New(level string) (*_logger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel // default
	}

	config := zap.Config{
		Encoding:         "json", // can be "console" for local dev
		Level:            zap.NewAtomicLevelAt(zapLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	zl, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &_logger{sugar: zl.Sugar()}, nil
}

func (l *_logger) Info(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

func (l *_logger) Infof(msg string, keysAndValues ...interface{}) {
	l.sugar.Infof(msg, keysAndValues...)
}

func (l *_logger) Error(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

func (l *_logger) Errorf(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorf(msg, keysAndValues...)
}

func (l *_logger) Debug(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

func (l *_logger) Warn(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

func (l *_logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

func (l *_logger) Sync() {
	_ = l.sugar.Sync()
}

// GetClient returns the singleton Redis client.
func SetLogger(logLevel string) {
	logg, _ := New(logLevel)
	Logger = logg
}
