package log

import "go.uber.org/zap"

var (
	sLogger LoggerInterface
)

// Register 注册
func Register(logger LoggerInterface) {
	if logger == nil {
		panic("logger interface is nil")
	}
	sLogger = logger
}

func With(args ...interface{}) *zap.SugaredLogger {
	return sLogger.With(args...)
}

func Sync() error {
	return sLogger.Sync()
}

func Debug(args ...interface{}) {
	sLogger.Debug(args...)
}

func Info(args ...interface{}) {
	sLogger.Info(args...)
}

func Warn(args ...interface{}) {
	sLogger.Warn(args...)
}

func Error(args ...interface{}) {
	sLogger.Error(args...)
}

func DPanic(args ...interface{}) {
	sLogger.DPanic(args...)
}

func Panic(args ...interface{}) {
	sLogger.Panic(args...)
}

func Fatal(args ...interface{}) {
	sLogger.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	sLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	sLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sLogger.Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	sLogger.DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	sLogger.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	sLogger.Fatalf(template, args...)
}
