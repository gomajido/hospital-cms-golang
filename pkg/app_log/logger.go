package app_log

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	stackdriver "github.com/tommy351/zap-stackdriver"
	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var zapLogger *zap.Logger

func Init() {
	_, file, line, _ := runtime.Caller(1)
	slash := strings.LastIndex(file, "/")
	file = file[slash+1:]
	if os.Getenv("APEXA_ENV") == "local" || os.Getenv("APEXA_ENV") == "staging" {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapLogger, _ = config.Build()
	} else {
		config := &zap.Config{
			Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
			Encoding:         "json",
			EncoderConfig:    stackdriver.EncoderConfig,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
		zapLogger, _ = config.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return &stackdriver.Core{
				Core: core,
			}
		}), zap.Fields(
			stackdriver.LogServiceContext(&stackdriver.ServiceContext{
				Service: "apexa-service",
				Version: "1.0.0",
			}),
		),
			zap.Fields(
				stackdriver.LogReportLocation(&stackdriver.ReportLocation{
					FilePath:     file,
					LineNumber:   line,
					FunctionName: "",
				}),
			),
		)
	}
}

func getLogger() *zap.Logger {
	if zapLogger == nil {
		Init()
	}
	return zapLogger
}

func Debug(msg string) {
	getLogger().Debug(msg)
}

func Debugf(format string, v ...any) {
	Debug(fmt.Sprintf(format, v...))
}

func Info(msg string) {
	getLogger().Info(msg)
}

func Infof(format string, v ...any) {
	Info(fmt.Sprintf(format, v...))
}

func Warn(msg string) {
	getLogger().Warn(msg)
}

func Warnf(format string, v ...any) {
	Warn(fmt.Sprintf(format, v...))
}

func Error(msg string) {
	getLogger().Error(msg)
}

func Errorf(format string, v ...any) {
	Error(fmt.Sprintf(format, v...))
}

func Fatal(msg string) {
	getLogger().Fatal(msg)
}

func Fatalf(format string, v ...any) {
	Fatal(fmt.Sprintf(format, v...))
}
