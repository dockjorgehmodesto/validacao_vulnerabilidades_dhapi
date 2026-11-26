package logging

import (
	"context"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var WrapperLogger *zap.Logger

func init() {
	// a fallback/root logger for events without context
	InitLoggers()
}

func InitLoggers() {
	level, err := zap.ParseAtomicLevel(viper.GetString("log.level"))
	if err != nil {
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	consoleEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		level,
	)
	Logger = zap.New(consoleCore, zap.AddCaller())
	WrapperLogger = Logger.WithOptions(zap.AddCallerSkip(2))
}

func getLogLevel() zapcore.Level {
	switch viper.GetString("log.level") {
	case "debug":
		return zapcore.DebugLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func LoggerWith(ctx context.Context) *zap.Logger {
	newLogger := Logger

	if ctx != nil {
		if reqID, ok := ctx.Value("request-id").(string); ok {
			newLogger = newLogger.With(zap.String("requestId", reqID))
		}
	}
	return newLogger
}
