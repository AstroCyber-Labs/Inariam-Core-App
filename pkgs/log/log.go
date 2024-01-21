package log

import (
	"github.com/fatih/color"
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

const (
	ErrInitializingLogger = "error initializing logger"
)

func init() {
	zapLogger, _ := zap.NewProduction()
	defer func(zapLogger *zap.Logger) {
		err := zapLogger.Sync()
		if err != nil {
			color.Red(err.Error())
		}
	}(zapLogger) // flushes buffer, if any
	Logger = zapLogger.Sugar()

	color.Green("[+] Logger is initialized")
}
