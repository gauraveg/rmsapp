package logger

import (
	"context"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type ZapLogger struct {
	Logger *zap.Logger
}

//type Wrapper struct {
//	logger Logger
//}

type Logger interface {
	DebugWithContext(context context.Context, args ...interface{})
	InfoWithContext(context context.Context, args ...interface{})
	WarnWithContext(context context.Context, args ...interface{})
	ErrorWithContext(context context.Context, args ...interface{})
	PanicWithContext(context context.Context, args ...interface{})
	FatalWithContext(context context.Context, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
}

func Init() *ZapLogger {
	config := zap.NewDevelopmentConfig()
	config.Encoding = "json"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "logCode"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.StacktraceKey = "stacktrace"

	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &ZapLogger{
		Logger: logger,
	}
}

func (l *ZapLogger) WithContext(ctx context.Context) *zap.Logger {
	requestId := ctx.Value("requestId").(string)
	childlog := l.Logger.With(zap.String("requestId", requestId))
	return childlog
}

func (l *ZapLogger) DebugWithContext(context context.Context, args ...interface{}) {
	l.WithContext(context)
	l.Debug(args...)
}

func (l *ZapLogger) Debug(args ...interface{}) {
	l.Logger.Debug(strconv.Itoa(int(zap.DebugLevel)), zap.Any("args", args))
}

func (l *ZapLogger) InfoWithContext(context context.Context, args ...interface{}) {
	childlog := l.WithContext(context)
	l.Info(childlog, args)
}

func (l *ZapLogger) Info(childlog *zap.Logger, args ...interface{}) {
	childlog.Info(strconv.Itoa(int(zap.InfoLevel)), zap.Any("args", args))
}

func (l *ZapLogger) WarnWithContext(context context.Context, args ...interface{}) {
	l.WithContext(context)
	l.Warn(args...)
}

func (l *ZapLogger) Warn(args ...interface{}) {
	l.Logger.Warn(strconv.Itoa(int(zap.WarnLevel)), zap.Any("args", args))
}

func (l *ZapLogger) ErrorWithContext(context context.Context, args ...interface{}) {
	childlog := l.WithContext(context)
	l.Error(childlog, args)
}

func (l *ZapLogger) Error(childlog *zap.Logger, args ...interface{}) {
	childlog.Error(strconv.Itoa(int(zap.ErrorLevel)), zap.Any("args", args))
}

func (l *ZapLogger) FatalWithContext(context context.Context, args ...interface{}) {
	l.WithContext(context)
	l.Fatal(args...)
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	l.Logger.Fatal(strconv.Itoa(int(zap.FatalLevel)), zap.Any("args", args))
}

// LogWrapperInit Initiate the logger and wrap
func LogWrapperInit() *ZapLogger {
	logger := Init()
	return logger
}

// GetLogContext - Get logger from the context
func GetLogContext(r *http.Request) *ZapLogger {
	logger, ok := r.Context().Value("logContext").(*ZapLogger)
	if !ok {
		return nil
	}
	return logger
}
