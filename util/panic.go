package util

import "go.uber.org/zap"

func Recover(log *zap.SugaredLogger) {
	if err := recover(); err != nil {
		log.Errorw("Recovered from panic", "error", err)
	}
}
