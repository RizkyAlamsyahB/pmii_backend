package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger initializes zap logger
func InitLogger(env string) error {
	var config zap.Config
	
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	
	var err error
	Log, err = config.Build()
	if err != nil {
		return err
	}
	
	return nil
}

// Info logs info message
func Info(message string, fields ...zap.Field) {
	Log.Info(message, fields...)
}

// Error logs error message
func Error(message string, fields ...zap.Field) {
	Log.Error(message, fields...)
}

// Debug logs debug message
func Debug(message string, fields ...zap.Field) {
	Log.Debug(message, fields...)
}

// Warn logs warning message
func Warn(message string, fields ...zap.Field) {
	Log.Warn(message, fields...)
}

// Fatal logs fatal message and exits
func Fatal(message string, fields ...zap.Field) {
	Log.Fatal(message, fields...)
}
