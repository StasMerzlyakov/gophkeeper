package app

import "github.com/sirupsen/logrus"

var log *logrus.Logger = logrus.New()

func SetMainLogger(lg *logrus.Logger) {
	log = lg
}

func GetMainLogger() *logrus.Logger {
	return log
}
