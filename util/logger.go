package util

import "go.uber.org/zap"

func InitLog(name string) *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugared := logger.Sugar()
	sugared = sugared.Named(name)
	return sugared
}
