package logger

import (
	"go.uber.org/zap"
)

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
