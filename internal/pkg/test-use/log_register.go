package test_use

import (
	"DatabaseManage/internal/pkg/log"
	"github.com/yuanbaopig/logger"
)

func BuildLogger() {
	logger.SetOptions(
		logger.WithLevel("debug"),
		logger.WithAddCallerSkip(1),
		logger.WithFormat("console"),
		logger.WithEnableColor(true),
	)
	log.Register(logger.Log.Sugar())
}
