package domain

import (
	"log"

	"go.uber.org/zap"
)

var mainLogger *zap.SugaredLogger

func GetApplicationLogger() *zap.SugaredLogger {

	if mainLogger != nil {
		return mainLogger
	} else {
		log.Default().Println("[WARN] application logger is not set")
		logger := zap.NewNop()
		mainLogger = logger.Sugar()
	}

	return mainLogger
}

func SetApplicationLogger(logger *zap.SugaredLogger) {
	mainLogger = logger
}
