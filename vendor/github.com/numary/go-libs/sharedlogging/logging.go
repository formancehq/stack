package sharedlogging

import "context"

type Logger interface {
	Debugf(fmt string, args ...interface{})
	Infof(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	WithFields(map[string]interface{}) Logger
	WithContext(ctx context.Context) Logger
}

type LoggerFactory interface {
	Get(ctx context.Context) Logger
}
type LoggerFactoryFn func(ctx context.Context) Logger

func (fn LoggerFactoryFn) Get(ctx context.Context) Logger {
	return fn(ctx)
}

func StaticLoggerFactory(l Logger) LoggerFactoryFn {
	return func(ctx context.Context) Logger {
		return l
	}
}

type noOpLogger struct{}

func (n noOpLogger) WithContext(ctx context.Context) Logger {
	return n
}

func (n noOpLogger) Debug(args ...interface{})              {}
func (n noOpLogger) Info(args ...interface{})               {}
func (n noOpLogger) Error(args ...interface{})              {}
func (n noOpLogger) Debugf(fmt string, args ...interface{}) {}
func (n noOpLogger) Infof(fmt string, args ...interface{})  {}
func (n noOpLogger) Errorf(fmt string, args ...interface{}) {}
func (n noOpLogger) WithFields(m map[string]interface{}) Logger {
	return n
}

var _ Logger = &noOpLogger{}

func NewNoOpLogger() *noOpLogger {
	return &noOpLogger{}
}

var loggerFactory LoggerFactory

func SetFactory(l LoggerFactory) {
	loggerFactory = l
}

func GetLogger(ctx context.Context) Logger {
	if loggerFactory == nil {
		return NewNoOpLogger()
	}
	return loggerFactory.Get(ctx)
}

func Debugf(fmt string, args ...interface{}) {
	GetLogger(context.Background()).Debugf(fmt, args...)
}
func Infof(fmt string, args ...interface{}) {
	GetLogger(context.Background()).Infof(fmt, args...)
}
func Errorf(fmt string, args ...interface{}) {
	GetLogger(context.Background()).Errorf(fmt, args...)
}
func Debug(args ...interface{}) {
	GetLogger(context.Background()).Debug(args...)
}
func Info(args ...interface{}) {
	GetLogger(context.Background()).Info(args...)
}
func Error(args ...interface{}) {
	GetLogger(context.Background()).Error(args...)
}
func WithFields(v map[string]interface{}) Logger {
	return GetLogger(context.Background()).WithFields(v)
}
