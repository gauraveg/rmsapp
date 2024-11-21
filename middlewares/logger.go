package middlewares

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type logCon string

const logContext logCon = "logContext"

func Init() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Encoding = "json"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.StacktraceKey = "stacktrace"

	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggers := Init()
		defer func(loggers *zap.Logger) {
			err := loggers.Sync()
			if err != nil {
				loggers.Error(err.Error())
				return
			}
		}(loggers)

		r = r.WithContext(context.WithValue(r.Context(), logContext, loggers))
		next.ServeHTTP(w, r)
	})
}

func LoggerContext(r *http.Request) *zap.Logger {
	logger, ok := r.Context().Value("logContext").(*zap.Logger)
	if !ok {
		return nil
	}
	return logger
}
