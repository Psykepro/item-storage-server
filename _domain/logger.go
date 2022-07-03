package _domain

type Logger interface {
	Debug(args ...any)
	Debugf(template string, args ...any)
	Info(args ...any)
	Infof(template string, args ...any)
	Warn(args ...any)
	Warnf(template string, args ...any)
	Error(args ...any)
	Errorf(template string, args ...any)
	DPanic(args ...any)
	DPanicf(template string, args ...any)
	Fatal(args ...any)
	Fatalf(template string, args ...any)
}
