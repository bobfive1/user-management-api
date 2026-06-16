package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	basedLogger *zap.SugaredLogger
	loggers     = make(map[string]*zap.SugaredLogger)
)

func GetDefaultLogger() *zap.SugaredLogger {
	logConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := logConfig.Build()
	return logger.Sugar()
}

func ApplyConfig(appName string, loglevel string) error {
	level, err := zap.ParseAtomicLevel(loglevel)
	if err != nil {
		return err
	}

	logConfig := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		InitialFields: map[string]interface{}{
			"application": appName,
		},
		Encoding:         "json",
		EncoderConfig:    NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := logConfig.Build()
	basedLogger = logger.Sugar()
	return nil
}

func GetLogger(name string) *zap.SugaredLogger {
	if logger, ok := loggers["name"]; ok {
		return logger
	}

	loggers[name] = basedLogger.Named(name)

	return basedLogger.Named(name)
}

func NewDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
