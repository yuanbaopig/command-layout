package log

import "go.uber.org/zap"

// LoggerInterface 定义了日志记录器应该具备的方法
type LoggerInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	DPanic(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Sync() error
	// 其他 zap.SugaredLogger 有的方法可以继续添加
	With(args ...interface{}) *zap.SugaredLogger
}
