package domain

import "go.uber.org/zap"

var mainLogger *zap.SugaredLogger

func GetApplicationLogger() *zap.SugaredLogger {
	return mainLogger
}

func SetApplicationLogger(logger *zap.SugaredLogger) {
	mainLogger = logger
}
